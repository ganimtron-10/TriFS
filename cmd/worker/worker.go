package main

import (
	"github.com/ganimtron-10/TriFS/internal/transport"
	"github.com/ganimtron-10/TriFS/internal/worker"
)

func StartWorker() {
	coreWorker, err := worker.CreateWorker()
	if err != nil {
		panic("Unable to create worker")
	}

	workerService := worker.CreateWorkerService(coreWorker)

	transport.StartRpcServer(coreWorker.Address, workerService)
}

func main() {
	StartWorker()
}
