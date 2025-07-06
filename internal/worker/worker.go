package worker

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/ganimtron-10/TriFS/internal/transport"
)

type WorkerConfig struct {
	MasterAddress     string
	Address           string
	HeartbeatInterval int
}

type Worker struct {
	*WorkerConfig
}

func getDefaultWorkerConfig() *WorkerConfig {
	return &WorkerConfig{
		MasterAddress:     common.DEFAULT_MASTER_ADDRESS,
		Address:           transport.GetAddressWithRandomPort(),
		HeartbeatInterval: 5,
	}
}

func (w *Worker) SendHeartBeat() {
	transport.DialRpcCall(w.MasterAddress, "MasterService.HeartBeat", &protocol.HeartBeatArgs{Address: w.Address}, &protocol.HeartBeatReply{})
}

func CreateWorker() (*Worker, error) {
	logger.Info(common.COMPONENT_WORKER, "Creating Worker...")

	worker := &Worker{
		getDefaultWorkerConfig(),
	}

	ticker := time.NewTicker(time.Second * time.Duration(worker.HeartbeatInterval))
	go func() {
		for {
			select {
			case <-ticker.C:
				worker.SendHeartBeat()
			}
		}
	}()

	if err := os.Mkdir(worker.Address, 0755); err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Unable to create directory named %s", worker.Address))
		return nil, fmt.Errorf("unable to initialize worker")
	}

	return worker, nil
}

func (w *Worker) AddConfig(config *WorkerConfig) *Worker {
	w.WorkerConfig = config
	return w
}

func (worker *Worker) handleWriteFile(filename string, data []byte) error {

	fullFilePath := path.Join(worker.Address, filename)
	if err := os.WriteFile(fullFilePath, data, 0755); err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Error while writing to file named %s. Error: %s", fullFilePath, err))
		return err
	}

	return nil
}
