package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

var commandList []*exec.Cmd

func spawn(spawnPath string) {
	command := exec.Command("go", "run", spawnPath)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	err := command.Start()
	if err != nil {
		logger.Error(common.COMPONENT_COMMON, fmt.Sprintf("Error while running %s", spawnPath))
	}

	commandList = append(commandList, command)
}

func main() {

	masterPath := "./cmd/master/main.go"
	workerPath := "./cmd/worker/main.go"
	clientPath := "./cmd/client/main.go"

	spawn(masterPath)
	time.Sleep(time.Second * 5)

	for i := 0; i < 3; i++ {
		spawn(workerPath)
		time.Sleep(time.Second * 1)
	}
	time.Sleep(time.Second * 5)

	spawn(clientPath)
	time.Sleep(time.Second * 5)

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	logger.Info(common.COMPONENT_COMMON, "Shutting down simulation...")

	for _, command := range commandList {
		if err := command.Process.Kill(); err != nil {
			fmt.Printf("Unable to kill Command %+v", command)
		}
	}

}
