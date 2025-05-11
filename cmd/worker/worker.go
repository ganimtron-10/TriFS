package main

import (
	"github.com/ganimtron-10/TriFS/internal/transport"
	"github.com/ganimtron-10/TriFS/internal/worker"
)

func StartWorker() {
	coreWorker := worker.CreateWorker()

	workerService := worker.CreateWorkerService(coreWorker)

	transport.StartRpcServer(coreWorker.Address, workerService)
}

func main() {
	StartWorker()
}
