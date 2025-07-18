package master

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MasterConfig struct {
	Port int
}

func getList(fileWorkerSet FileWorkerSet) []string {
	workerList := make([]string, 0, len(fileWorkerSet))
	for worker := range fileWorkerSet {
		workerList = append(workerList, worker)
	}
	return workerList
}

type WorkerInfo struct {
	FileHashes map[string]struct{}
}

type FileWorkerSet = map[string]bool

type Master struct {
	protocol.UnimplementedMasterServiceServer

	*MasterConfig
	WorkerPool            map[string]*WorkerInfo
	WorkerPoolLock        sync.RWMutex
	FileHashWorkerMap     map[string]FileWorkerSet
	FileHashWorkerMapLock sync.RWMutex
}

func getDefaultMasterConfig() *MasterConfig {
	return &MasterConfig{
		Port: common.DEFAULT_MASTER_PORT,
	}
}

func createMaster() *Master {
	logger.Info(common.COMPONENT_MASTER, "Creating Master...")
	return &Master{
		MasterConfig:      getDefaultMasterConfig(),
		WorkerPool:        make(map[string]*WorkerInfo),
		FileHashWorkerMap: make(map[string]FileWorkerSet),
	}
}

func (m *Master) AddConfig(config *MasterConfig) *Master {
	m.MasterConfig = config
	return m
}

func StartMaster() error {
	master := createMaster()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", master.Port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	protocol.RegisterMasterServiceServer(grpcServer, master)

	go (func() {
		logger.Info(common.COMPONENT_MASTER, "Starting Master grpcServer...")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error(common.COMPONENT_MASTER, fmt.Sprintf("Master grpcServer failed to server: %+v", err))
		}
	})()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logger.Info(common.COMPONENT_MASTER, "Shutting down Master grpcServer...")

	grpcServer.GracefulStop()
	logger.Info(common.COMPONENT_MASTER, "Stopped Master grpcServer")

	return nil
}

func (m *Master) GetFileWorkers(ctx context.Context, req *protocol.GetFileWorkersRequest) (*protocol.GetFileWorkersResponse, error) {
	if err := common.ValidateRequest(req); err != nil {
		return nil, err
	}

	workerUrls, err := m.handleGetFileWorkers(req.Filename)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get workers for file %s: %+v", req.Filename, err)
	}

	res := &protocol.GetFileWorkersResponse{
		WorkerUrls: workerUrls,
	}

	return res, nil
}

func (m *Master) AllocateFileWorkers(ctx context.Context, req *protocol.AllocateFileWorkersRequest) (*protocol.AllocateFileWorkersResponse, error) {
	if err := common.ValidateRequest(req); err != nil {
		return nil, err
	}

	workerUrls, err := m.handleAllocateFileWorkers(req.Filename)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to allocate workers for file %s: %+v", req.Filename, err)
	}

	res := &protocol.AllocateFileWorkersResponse{WorkerUrls: workerUrls}

	return res, nil
}

func (m *Master) Heartbeat(ctx context.Context, req *protocol.HeartbeatRequest) (*protocol.HeartbeatResponse, error) {
	if err := common.ValidateRequest(req); err != nil {
		return nil, err
	}

	// TODO: Check later whether function needs to return anything
	fileHashes := make(map[string]struct{})
	for _, fileHash := range req.HostedFileHashes {
		fileHashes[fileHash] = struct{}{}
	}

	m.handleHeartbeat(req.WorkerAddress, fileHashes)

	return &protocol.HeartbeatResponse{}, nil
}
