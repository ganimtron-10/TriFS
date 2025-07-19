package worker

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type WorkerConfig struct {
	MasterAddress     string
	Address           string
	HeartbeatInterval int
}

type FileInfo struct {
	PackId string
	Offset int
	Size   int
}

type Worker struct {
	protocol.UnimplementedWorkerServiceServer

	*WorkerConfig
	fileStore     map[string]*FileInfo
	fileStoreLock sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func getDefaultWorkerConfig() *WorkerConfig {
	return &WorkerConfig{
		MasterAddress:     common.DEFAULT_MASTER_ADDRESS,
		Address:           common.GetAddressWithRandomPort(),
		HeartbeatInterval: 5,
	}
}

func getFileHashes(fileStore map[string]*FileInfo) []string {
	fileHashes := []string{}
	for fileHash := range fileStore {
		fileHashes = append(fileHashes, fileHash)
	}
	return fileHashes
}

func (w *Worker) SendHeartBeat(masterClient protocol.MasterServiceClient) {
	w.fileStoreLock.RLock()
	defer w.fileStoreLock.RUnlock()

	ctx, cancel := context.WithTimeout(w.ctx, time.Second)
	defer cancel()
	_, err := masterClient.Heartbeat(ctx, &protocol.HeartbeatRequest{
		WorkerAddress:    w.Address,
		HostedFileHashes: getFileHashes(w.fileStore),
	})
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Sending Heartbeat failed", "reason", err.Error())
	}
	// logger.Info(common.COMPONENT_WORKER, "Heartbeat sent successfully")
}

func createWorker() (*Worker, error) {
	logger.Info(common.COMPONENT_WORKER, "Creating Worker...")

	ctx, cancel := context.WithCancel(context.Background())

	worker := &Worker{
		WorkerConfig: getDefaultWorkerConfig(),
		fileStore:    make(map[string]*FileInfo),
		ctx:          ctx,
		cancel:       cancel,
	}

	if err := os.Mkdir(worker.Address, 0755); err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Unable to create directory named %s", worker.Address))
		return nil, fmt.Errorf("unable to initialize worker")
	}

	return worker, nil
}

func (w *Worker) Shutdown() {
	w.cancel()
	w.wg.Wait()
}

func (w *Worker) AddConfig(config *WorkerConfig) *Worker {
	w.WorkerConfig = config
	return w
}

func (w *Worker) startHeartbeating(ctx context.Context) error {
	defer w.wg.Done()

	grpcClient, err := grpc.NewClient(w.MasterAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("unable to create grpc client for sending heartbeat: %+v", err)
	}
	defer grpcClient.Close()
	masterClient := protocol.NewMasterServiceClient(grpcClient)

	ticker := time.NewTicker(time.Second * time.Duration(w.HeartbeatInterval))
	defer ticker.Stop()

	logger.Info(common.COMPONENT_WORKER, "Starting Heartbeating...")
	for {
		select {
		case <-ctx.Done():
			logger.Info(common.COMPONENT_WORKER, "Stopping Heartbeating", "reason", ctx.Err().Error())
			return nil
		case <-ticker.C:
			w.SendHeartBeat(masterClient)
		}
	}
}

func StartWorker() error {
	worker, err := createWorker()
	if err != nil {
		return fmt.Errorf("unable to create worker: %+v", err)
	}

	lis, err := net.Listen("tcp", worker.Address)
	if err != nil {
		return fmt.Errorf("unable to create listner: %+v", err)
	}

	grpcServer := grpc.NewServer()
	protocol.RegisterWorkerServiceServer(grpcServer, worker)

	worker.wg.Add(1)
	go func() {
		defer worker.wg.Done()
		logger.Info(common.COMPONENT_WORKER, "Starting Worker grpcServer...")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Worker grpcServer failed to serve: %+v", err))
		}
	}()

	worker.wg.Add(1)
	go worker.startHeartbeating(worker.ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logger.Info(common.COMPONENT_WORKER, "Shutting down Worker grpcServer...")

	grpcServer.GracefulStop()
	worker.Shutdown()
	logger.Info(common.COMPONENT_WORKER, "Stopped Worker grpcServer")

	return nil
}

func (w *Worker) ReadFile(ctx context.Context, req *protocol.ReadFileRequest) (*protocol.ReadFileResponse, error) {
	if err := common.ValidateRequest(req); err != nil {
		return nil, err
	}

	fileData, err := w.handleReadFile(req.Filename)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Error while handling ReadFile. Error: %s", err))
		return nil, status.Errorf(codes.Internal, "unable to read file %s: %+v", req.Filename, err)
	}

	res := &protocol.ReadFileResponse{
		Filename: req.Filename,
		Data:     fileData,
	}

	return res, nil

}
func (w *Worker) WriteFile(ctx context.Context, req *protocol.WriteFileRequest) (*protocol.WriteFileResponse, error) {
	if err := common.ValidateRequest(req); err != nil {
		return nil, err
	}

	err := w.handleWriteFile(req.Filename, req.Data)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Error while handling WriteFile. Error: %s", err))
		return nil, status.Errorf(codes.Internal, "unable to write file %s: %+v", req.Filename, err)
	}

	return &protocol.WriteFileResponse{}, nil

}
