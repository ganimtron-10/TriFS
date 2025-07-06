package master

import (
	"fmt"
	"sync"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

type MasterConfig struct {
	Port int
}

type Master struct {
	*MasterConfig
	WorkerPool     map[string]int
	WorkerPoolLock sync.RWMutex
}

func getDefaultMasterConfig() *MasterConfig {
	return &MasterConfig{
		Port: common.DEFAULT_MASTER_PORT,
	}
}

func CreateMaster() *Master {
	logger.Info(common.COMPONENT_MASTER, "Creating Master...")
	return &Master{
		MasterConfig: getDefaultMasterConfig(),
		WorkerPool:   make(map[string]int),
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
	// choose and return the worker url to write to the file

	master.WorkerPoolLock.RLock()
	for key := range master.WorkerPool {
		return []byte(key), nil
	}
	master.WorkerPoolLock.RUnlock()

	logger.Error(common.COMPONENT_MASTER, "No Worker URL in WorkerPool")
	return nil, fmt.Errorf("worker not available. please try later")
}
