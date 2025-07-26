package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

var commandList []*exec.Cmd
var cwd string

func spawn(spawnPath string) {
	command := exec.Command("go", "run", spawnPath)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Dir = cwd

	err := command.Start()
	if err != nil {
		logger.Error(common.COMPONENT_COMMON, fmt.Sprintf("Error while running %s", spawnPath))
	}

	commandList = append(commandList, command)
}

func main() {

	cwd = ".simulate"
	err := os.MkdirAll(cwd, 0644)
	if err != nil {
		logger.Error(common.COMPONENT_COMMON, fmt.Sprintf("Unable to create directory named %s", cwd))
		log.Fatalf("Unable to create directory named %s", cwd)
	}

	defer func() {
		err := os.RemoveAll("./.simulate")
		if err != nil {
			logger.Error(common.COMPONENT_COMMON, fmt.Sprintf("Unable to delete directory named %s", cwd))
			log.Fatalf("Unable to delete directory named %s", cwd)
		}
	}()

	timeInterval := time.Second * 5

	masterPath := "../cmd/master/main.go"
	workerPath := "../cmd/worker/main.go"
	clientPath := "../cmd/client/main.go"

	spawn(masterPath)
	time.Sleep(timeInterval)

	for i := 0; i < 3; i++ {
		spawn(workerPath)
		time.Sleep(time.Second * 1)
	}
	time.Sleep(timeInterval)

	spawn(clientPath)
	time.Sleep(timeInterval)

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Wait for everything to complete
	time.Sleep(5 * time.Second)

	logger.Info(common.COMPONENT_COMMON, "Shutting down simulation...")
	for _, command := range commandList {
		if err := command.Process.Kill(); err != nil {
			fmt.Printf("Unable to kill Command %+v", command)
		}
	}

}
