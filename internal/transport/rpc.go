package transport

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func SendRpcCall(address string, rpcServiceName string, rpcRequest any, rpcResponse any) error {
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

func StartRpcServer(port int) (*rpc.Client, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error creating listner", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}

func RegisterService(service any) {
	rpc.Register(service)
}
