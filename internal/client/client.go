package client

import (
	"fmt"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/ganimtron-10/TriFS/internal/transport"
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
		getDefaultClientConfig(),
	}
}

func (client *Client) AddConfig(config *ClientConfig) *Client {
	client.ClientConfig = config
	return client
}

func (client *Client) Read(filename string) error {
	requestArgs := &protocol.ReadFileRequestArgs{Filename: filename}
	requestReply := &protocol.ReadFileRequestReply{}

	err := transport.DialRpcCall(client.MasterAddress, "MasterService.ReadFile", requestArgs, requestReply)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Master ReadFile Error: %s", err))
		return err
	}

	logger.Info(common.COMPONENT_CLIENT, "Master ReadFile Response", "WorkerUrls", requestReply.WorkerUrls)

	args := &protocol.ReadFileArgs{Filename: filename}
	reply := &protocol.ReadFileReply{}

	err = transport.DialRpcCall(requestReply.WorkerUrls[0], "WorkerService.ReadFile", args, reply)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Worker ReadFile Error: %s", err))
		return err
	}

	logger.Info(common.COMPONENT_CLIENT, "Worker ReadFile Response", "Data", string(reply.Data))
	return nil
}

func (client *Client) Write(filename, data string) {
	requestArgs := &protocol.WriteFileRequestArgs{Filename: filename}
	requestReply := &protocol.WriteFileRequestReply{}

	err := transport.DialRpcCall(client.MasterAddress, "MasterService.WriteFile", requestArgs, requestReply)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Master WriteFile Error: %s", err))
	}

	logger.Info(common.COMPONENT_CLIENT, "Master WriteFile Response", "WorkerUrl", requestReply.WorkerUrls)

	args := &protocol.WriteFileArgs{Filename: filename, Data: []byte(data)}
	reply := &protocol.WriteFileReply{}

	err = transport.DialRpcCall(requestReply.WorkerUrls[0], "WorkerService.WriteFile", args, reply)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Worker WriteFile Error: %s", err))
	}

	logger.Info(common.COMPONENT_CLIENT, "Worker Successfully Wrote File")
}
