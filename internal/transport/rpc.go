package transport

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

func DialRpcCall(address string, rpcServiceName string, rpcRequest any, rpcResponse any) error {
	logger.Info(common.COMPONENT_COMMON, fmt.Sprintf("Dialing RPC Call to %s", rpcServiceName))
	rpcClient, err := rpc.Dial("tcp", address)
	if err != nil {
		return err
	}

	err = rpcClient.Call(rpcServiceName, rpcRequest, rpcResponse)
	if err != nil {
		return err
	}

	return nil
}

func RegisterServices(services []interface{}) {
	logger.Info(common.COMPONENT_COMMON, "Registering Services...")
	for _, service := range services {
		rpc.Register(service)
	}
}

func StartRpcServer(address string, services ...interface{}) (*rpc.Client, error) {

	RegisterServices(services)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error(common.COMPONENT_COMMON, fmt.Sprint("Error creating listener", err))
	}
	defer listener.Close()

	logger.Info(common.COMPONENT_COMMON, fmt.Sprintf("Accepting Connections on %s", listener.Addr().String()))
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error(common.COMPONENT_COMMON, fmt.Sprint("Error accepting connection", err))
			continue
		}

		go rpc.ServeConn(conn)
	}
}

func GetAddressWithRandomPort() string {
	listener, err := net.Listen("tcp", "0")
	if err != nil {
		logger.Error(common.COMPONENT_COMMON, fmt.Sprint("Error creating listner", err))
	}
	defer listener.Close()

	return listener.Addr().String()
}
