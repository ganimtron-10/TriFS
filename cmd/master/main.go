package main

import (
	"github.com/ganimtron-10/TriFS/internal/master"
	"github.com/ganimtron-10/TriFS/internal/transport"
)

func StartMaster() {
	coreMaster := master.CreateMaster()

	masterService := master.CreateMasterService(coreMaster)

	transport.StartRpcServer(coreMaster.Port, masterService)
}

func main() {
	StartMaster()
}
