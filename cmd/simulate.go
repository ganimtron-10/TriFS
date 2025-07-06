package main

import (
	"fmt"
	"time"

	"github.com/ganimtron-10/TriFS/internal/client"
	"github.com/ganimtron-10/TriFS/internal/master"
	"github.com/ganimtron-10/TriFS/internal/transport"
	"github.com/ganimtron-10/TriFS/internal/worker"
)

func StartMaster() {
	coreMaster := master.CreateMaster()

	masterService := master.CreateMasterService(coreMaster)

	transport.StartRpcServer(fmt.Sprintf(":%d", coreMaster.Port), masterService)

}

func StartWorker() {
	coreWorker := worker.CreateWorker()

	workerService := worker.CreateWorkerService(coreWorker)

	transport.StartRpcServer(coreWorker.Address, workerService)

}

func main() {

	go StartMaster()

	go StartWorker()
	go StartWorker()
	go StartWorker()

	time.Sleep(time.Second * 20)

	tc := client.CreateClient()
	tc.Read("test.txt")

	select {}
}
