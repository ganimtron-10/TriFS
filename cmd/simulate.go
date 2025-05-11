package main

import (
	"time"

	"github.com/ganimtron-10/TriFS/internal/client"
	"github.com/ganimtron-10/TriFS/internal/master"
	"github.com/ganimtron-10/TriFS/internal/transport"
)

func StartMaster() {
	coreMaster := master.CreateMaster()

	masterService := master.CreateMasterService(coreMaster)

	transport.StartRpcServer(coreMaster.Port, masterService)

}

func main() {

	go StartMaster()

	time.Sleep(time.Second * 5)

	tc := client.CreateClient()
	tc.Read("test.txt")
}
