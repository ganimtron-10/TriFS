package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganimtron-10/TriFS/internal/client"
	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
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

	tc.Write("test1.txt", "Test File 1")
	time.Sleep(time.Second * 5)
	tc.Write("test2.txt", "Test File 2")
	time.Sleep(time.Second * 5)
	tc.Write("test3.txt", "Test File 3")
	time.Sleep(time.Second * 5)
	tc.Write("test4.txt", "Test File 4")
	time.Sleep(time.Second * 5)

	tc.Read("test1.txt")
	time.Sleep(time.Second * 5)
	tc.Read("test3.txt")
	time.Sleep(time.Second * 5)

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	logger.Info(common.COMPONENT_COMMON, "Shutting down simulation...")
}
