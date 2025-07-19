package client

import (
	"context"
	"fmt"
	"time"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	MasterAddress string
}

type Client struct {
	*ClientConfig
}

func getDefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		MasterAddress: common.DEFAULT_MASTER_ADDRESS,
	}
}

func CreateClient() *Client {
	logger.Info(common.COMPONENT_CLIENT, "Creating Client...")
	return &Client{
		ClientConfig: getDefaultClientConfig(),
	}
}

func (c *Client) AddConfig(config *ClientConfig) *Client {
	c.ClientConfig = config
	return c
}

func dialGRPC(address string) (*grpc.ClientConn, error) {
	return grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (c *Client) Read(filename string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conn, err := dialGRPC(c.MasterAddress)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Failed to connect to Master: %s", err))
		return err
	}
	defer conn.Close()

	masterClient := protocol.NewMasterServiceClient(conn)
	req := &protocol.GetFileWorkersRequest{Filename: filename}
	res, err := masterClient.GetFileWorkers(ctx, req)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Master GetFileWorkers error: %s", err))
		return err
	}

	if len(res.WorkerUrls) == 0 {
		return fmt.Errorf("no worker urls returned from master")
	}

	logger.Info(common.COMPONENT_CLIENT, "Master GetFileWorkers Response", "WorkerUrls", res.WorkerUrls)

	// TODO: Add Retry and Fallback Logic
	workerConn, err := dialGRPC(res.WorkerUrls[0])
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Failed to connect to Worker: %s", err))
		return err
	}
	defer workerConn.Close()

	workerClient := protocol.NewWorkerServiceClient(workerConn)
	wreq := &protocol.ReadFileRequest{Filename: filename}
	wres, err := workerClient.ReadFile(ctx, wreq)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Worker ReadFile error: %s", err))
		return err
	}

	logger.Info(common.COMPONENT_CLIENT, "Worker ReadFile Response", "Data", string(wres.Data))
	return nil
}

func (c *Client) Write(filename, data string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conn, err := dialGRPC(c.MasterAddress)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Failed to connect to Master: %s", err))
		return err
	}
	defer conn.Close()

	masterClient := protocol.NewMasterServiceClient(conn)
	req := &protocol.AllocateFileWorkersRequest{Filename: filename}
	res, err := masterClient.AllocateFileWorkers(ctx, req)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Master AllocateFileWorkers error: %s", err))
		return err
	}

	if len(res.WorkerUrls) == 0 {
		return fmt.Errorf("no worker urls returned from master")
	}

	logger.Info(common.COMPONENT_CLIENT, "Master AllocateFileWorkers Response", "WorkerUrls", res.WorkerUrls)

	// TODO: Add Retry and Fallback Logic
	workerConn, err := dialGRPC(res.WorkerUrls[0])
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Failed to connect to Worker: %s", err))
		return err
	}
	defer workerConn.Close()

	workerClient := protocol.NewWorkerServiceClient(workerConn)
	wreq := &protocol.WriteFileRequest{
		Filename: filename,
		Data:     []byte(data),
	}
	_, err = workerClient.WriteFile(ctx, wreq)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Worker WriteFile error: %s", err))
		return err
	}

	logger.Info(common.COMPONENT_CLIENT, "Worker Successfully Wrote File")
	return nil
}
