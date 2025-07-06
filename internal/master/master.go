package master

import (
	"fmt"
	"sync"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

type MasterConfig struct {
	Port           int
	WorkerPool     map[string]int
	WorkerPoolLock sync.RWMutex
}

type Master struct {
	*MasterConfig
}

func getDefaultMasterConfig() *MasterConfig {
	return &MasterConfig{
		Port:       common.DEFAULT_MASTER_PORT,
		WorkerPool: make(map[string]int),
	}
}

func CreateMaster() *Master {
	logger.Info(common.COMPONENT_MASTER, "Creating Master...")
	return &Master{
		getDefaultMasterConfig(),
	}
}

func (master *Master) AddConfig(config *MasterConfig) *Master {
	master.MasterConfig = config
	return master
}

func (master *Master) handleReadFile(filename string) ([]byte, error) {
	// return the worker url to access the file
	fmt.Println(master.WorkerPool)

	return []byte{0, 1, 2, 3, 4, 5}, nil
}

func (master *Master) handleWriteFileRequest(filename string) ([]byte, error) {
	// return the worker url to write to the file

	for key, _ := range master.WorkerPool {
		return []byte(key), nil
	}

	logger.Error(common.COMPONENT_MASTER, "Error while handling WriteFile: ", "No Workers Available")
	return nil, fmt.Errorf("No Workers Available")
}
