package main

import (
	"log"

	"github.com/ganimtron-10/TriFS/internal/worker"
)

func main() {
	if err := worker.StartWorker(); err != nil {
		log.Fatalf("Unable to Start Worker Server: %+v", err)
	}
}
