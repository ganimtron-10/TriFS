package master

import (
	"log"

	"github.com/ganimtron-10/TriFS/internal/service"
	"github.com/ganimtron-10/TriFS/internal/transport"
)

type MasterConfig struct {
	Port int
}

type Master struct {
	*MasterConfig
}

func getDefaultMasterConfig() *MasterConfig {
	return &MasterConfig{
		Port: 9867,
	}
}

func CreateMaster() *Master {
	return &Master{
		getDefaultMasterConfig(),
	}
}

func (master *Master) AddConfig(config *MasterConfig) *Master {
	master.MasterConfig = config
	return master
}

func RegisterAllService() {
	log.Println("Registering Services...")
	transport.RegisterService(new(service.FileService))
}

func (master *Master) Start() {
	RegisterAllService()
	transport.StartRpcServer(master.Port)
}

func Init() {
	master := CreateMaster()
	log.Println("Master Initializing...")
	master.Start()
}
