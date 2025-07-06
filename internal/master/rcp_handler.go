package master

import (
	"fmt"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
)

type MasterService struct {
	master *Master
}

func CreateMasterService(master *Master) *MasterService {
	logger.Info(common.COMPONENT_MASTER, "Creating Master Service...")
	return &MasterService{
		master: master,
	}
}

func validateArgsNReply(args, reply any) error {
	if args == nil {
		return fmt.Errorf("RPC Args is empty")
	}
	if reply == nil {
		return fmt.Errorf("RPC Reply is empty")
	}

	return nil
}

func (s *MasterService) ReadFile(args *protocol.ReadFileArgs, reply *protocol.ReadFileReply) error {
	if err := validateArgsNReply(args, reply); err != nil {
		return err
	}

	fileData, err := s.master.handleReadFile(args.Filename)
	if err != nil {
		return err
	}

	reply.Filename = args.Filename
	reply.Data = fileData

	return nil
}

func (s *MasterService) HeartBeat(args *protocol.HeartBeatArgs, reply *protocol.HeartBeatReply) error {
	if err := validateArgsNReply(args, reply); err != nil {
		return err
	}

	s.master.WorkerPoolLock.Lock()
	s.master.WorkerPool[args.Address] = 1
	s.master.WorkerPoolLock.Unlock()

	return nil
}
