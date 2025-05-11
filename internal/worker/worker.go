package worker

import (
	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/ganimtron-10/TriFS/internal/transport"
)

type WorkerConfig struct {
	MasterAddress string
	Address       string
}

type Worker struct {
	*WorkerConfig
}

func getDefaultWorkerConfig() *WorkerConfig {
	return &WorkerConfig{
		MasterAddress: common.DEFAULT_MASTER_ADDRESS,
		Address:       transport.GetAddressWithRandomPort(),
	}
}

func (w *Worker) SendHeartBeat() {
	transport.DialRpcCall(w.MasterAddress, "MasterService.HeartBeat", &protocol.HeartBeatArgs{Address: w.Address}, &protocol.HeartBeatReply{})
}

func CreateWorker() *Worker {
	logger.Info(common.COMPONENT_WORKER, "Creating Worker...")
	return &Worker{
		getDefaultWorkerConfig(),
	}
}

func (w *Worker) AddConfig(config *WorkerConfig) *Worker {
	w.WorkerConfig = config
	return w
}
