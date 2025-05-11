package client

import (
	"log"

	"github.com/ganimtron-10/TriFS/internal/common"
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
	log.Println("Creating Client...")
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
		log.Println("Error while performing ReadFile", err)
	}

	log.Println("Response Data: ", reply.Data)
}
