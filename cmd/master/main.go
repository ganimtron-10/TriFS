package main

import (
	"log"

	"github.com/ganimtron-10/TriFS/internal/master"
)

func main() {
	if err := master.StartMaster(); err != nil {
		log.Fatalf("Unable to Start Master Server: %+v", err)
	}
}
