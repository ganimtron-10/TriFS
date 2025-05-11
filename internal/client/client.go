package client

import (
	"log"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/service"
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
		MasterAddress: "localhost:9867",
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

func CreateReadFileRequest(filename string) *service.ReadFileRequest {
	return &service.ReadFileRequest{
		Message: &service.Message{
			Code: common.MESSAGE_READ,
		},
		Filename: filename,
	}
}

func (client *Client) Read(filename string) {
	response := &service.ReadFileResponse{}

	err := transport.SendRpcCall(client.MasterAddress, "FileService.Read", CreateReadFileRequest(filename), response)
	if err != nil {
		log.Println("Error while performing File Read", err)
	}

	log.Println("Response Data: ", response.Data)
}
