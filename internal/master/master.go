package master

import (
	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

type MasterConfig struct {
	Port int
}

type Master struct {
	*MasterConfig
}

func getDefaultMasterConfig() *MasterConfig {
	return &MasterConfig{
		Port: common.DEFAULT_MASTER_PORT,
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
	return []byte{0, 1, 2, 3, 4, 5}, nil
}
