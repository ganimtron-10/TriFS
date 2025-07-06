package worker

import (
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

func CreateWorker() *Worker {
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

	return worker
}

func (w *Worker) AddConfig(config *WorkerConfig) *Worker {
	w.WorkerConfig = config
	return w
}
