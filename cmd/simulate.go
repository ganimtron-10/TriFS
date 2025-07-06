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
	coreWorker, err := worker.CreateWorker()
	if err != nil {
		panic("Unable to create worker")
	}

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
	tc.Write("test.txt", "This is the data that is to be written to the file")

	select {}
}
