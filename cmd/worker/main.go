package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ganimtron-10/TriFS/internal/transport"
	"github.com/ganimtron-10/TriFS/internal/worker"
)

func StartWorker() {
	coreWorker, err := worker.CreateWorker()
	if err != nil {
		panic("Unable to create worker")
	}

	workerService := worker.CreateWorkerService(coreWorker)

	go func() {
		if err := transport.StartRpcServer(coreWorker.Address, workerService); err != nil {
			log.Fatalf("RPC server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Worker shutting down...")
}

func main() {
	StartWorker()
}
