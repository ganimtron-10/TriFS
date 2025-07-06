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

func (client *Client) Read(filename string) {
	args := &protocol.ReadFileArgs{Filename: filename}
	reply := &protocol.ReadFileReply{}

	err := transport.DialRpcCall(client.MasterAddress, "MasterService.ReadFile", args, reply)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Master ReadFile Error: %s", err))
	}

	logger.Info(common.COMPONENT_CLIENT, "Master ReadFile Response", "Data", reply.Data)
}

func (client *Client) Write(filename, data string) {
	requestArgs := &protocol.WriteFileRequestArgs{Filename: filename}
	requestReply := &protocol.WriteFileRequestReply{}

	err := transport.DialRpcCall(client.MasterAddress, "MasterService.WriteFile", requestArgs, requestReply)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Master WriteFile Error: %s", err))
	}

	logger.Info(common.COMPONENT_CLIENT, "Master WriteFile Response", "WorkerUrl", requestReply.WorkerUrl)

	args := &protocol.WriteFileArgs{Filename: filename, Data: []byte(data)}
	reply := &protocol.WriteFileReply{}

	err = transport.DialRpcCall(client.MasterAddress, "WorkerService.WriteFile", args, reply)
	if err != nil {
		logger.Error(common.COMPONENT_CLIENT, fmt.Sprintf("Worker WriteFile Error: %s", err))
	}

	logger.Info(common.COMPONENT_CLIENT, "Worker Successfully Wrote File")
}
