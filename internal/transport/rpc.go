package transport

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func DialRpcCall(address string, rpcServiceName string, rpcRequest any, rpcResponse any) error {
	log.Printf("Dialing RPC Call to %s", rpcServiceName)
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
	log.Println("Registering Services...")
	for _, service := range services {
		rpc.Register(service)
	}
}

func StartRpcServer(port int, services ...interface{}) (*rpc.Client, error) {

	RegisterServices(services)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error creating listner", err)
	}
	defer listener.Close()

	log.Printf("Accepting Connections on Port :%d", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}
