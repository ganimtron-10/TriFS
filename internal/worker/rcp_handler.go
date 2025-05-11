package worker

import (
	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

type WorkerService struct {
	worker *Worker
}

func CreateWorkerService(worker *Worker) *WorkerService {
	logger.Info(common.COMPONENT_WORKER, "Creating Worker Service...")
	return &WorkerService{
		worker: worker,
	}
}
