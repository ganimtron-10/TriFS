package worker

import (
	"fmt"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
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

func (s *WorkerService) WriteFile(args *protocol.WriteFileArgs, reply *protocol.WriteFileReply) error {
	if err := common.ValidateArgsNReply(args, reply); err != nil {
		return err
	}

	err := s.worker.handleWriteFile(args.Filename, args.Data)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Error while handling WriteFile %s", err))
		return err
	}

	return nil
}
