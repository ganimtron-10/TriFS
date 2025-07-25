package worker

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type WorkerConfig struct {
	MasterAddress     string
	Address           string
	HeartbeatInterval int
	Id                string
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
	WAL           WAL
}

func getDefaultWorkerConfig() *WorkerConfig {
	return &WorkerConfig{
		MasterAddress:     common.DEFAULT_MASTER_ADDRESS,
		Address:           common.GetAddressWithRandomPort(),
		HeartbeatInterval: 5,
		Id:                uuid.NewString()[:8],
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

	defaultWorkerConfig := getDefaultWorkerConfig()
	worker := &Worker{
		WorkerConfig: defaultWorkerConfig,
		fileStore:    make(map[string]*FileInfo),
		WAL:          createWAL(defaultWorkerConfig.Id),
		ctx:          ctx,
		cancel:       cancel,
	}

	worker.createWorkerDirectoryStructure()

	return worker, nil
}

func (w *Worker) Shutdown() {
	_, err := w.WAL.flushToFile()
	if err != nil {
		logger.Info(common.COMPONENT_WORKER, "Unable to flush WAL", "error", err)
	}

	w.cancel()
	w.wg.Wait()
}

func (w *Worker) createWorkerDirectoryStructure() error {
	baseDir := w.Id
	subDirs := []string{
		common.FOLDER_DATA,
		common.FOLDER_WAL,
		common.FOLDER_PACK,
	}

	if err := os.MkdirAll(baseDir, 0644); err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Unable to create base worker directory %s: %v", baseDir, err))
		return fmt.Errorf("unable to create worker dirs: %w", err)
	}
	logger.Info(common.COMPONENT_WORKER, fmt.Sprintf("Worker base directory %s created or already exists.", baseDir))

	for _, subDir := range subDirs {

		fullPath := filepath.Join(baseDir, subDir)

		if err := os.MkdirAll(fullPath, 0644); err != nil {
			logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Unable to create subdirectory %s: %v", fullPath, err))
			return fmt.Errorf("unable to create worker dirs: %w", err)
		}
		logger.Info(common.COMPONENT_WORKER, fmt.Sprintf("Worker subdirectory %s created or already exists.", fullPath))
	}

	return nil
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
	// TODO: Make sure Master is Up, due to lazy grpc conn it doesnt throw error if not connected

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
	go func() {
		if err := worker.startHeartbeating(worker.ctx); err != nil {
			logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Heartbeating failed: %v", err))
		}
	}()

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
