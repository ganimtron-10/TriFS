package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganimtron-10/TriFS/internal/client"
	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

func main() {

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
