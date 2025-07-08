package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

func spawn(spawnPath string) {
	command := exec.Command("go", "run", spawnPath)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	err := command.Start()
	if err != nil {
		logger.Error(common.COMPONENT_COMMON, fmt.Sprintf("Error while running %s", spawnPath))
	}
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

}
