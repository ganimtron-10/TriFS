package main

import (
	"fmt"

	"github.com/ganimtron-10/TriFS/internal/master"
	"github.com/ganimtron-10/TriFS/internal/transport"
)

func StartMaster() {
	coreMaster := master.CreateMaster()

	masterService := master.CreateMasterService(coreMaster)

	transport.StartRpcServer(fmt.Sprintf(":%d", coreMaster.Port), masterService)
}

func main() {
	StartMaster()
}
