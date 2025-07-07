package master

import (
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

func (s *MasterService) ReadFile(args *protocol.ReadFileRequestArgs, reply *protocol.ReadFileRequestReply) error {
	if err := common.ValidateArgsNReply(args, reply); err != nil {
		return err
	}

	workerUrls, err := s.master.handleReadFile(args.Filename)
	if err != nil {
		return err
	}

	reply.WorkerUrls = workerUrls

	return nil
}

func (s *MasterService) HeartBeat(args *protocol.HeartBeatArgs, reply *protocol.HeartBeatReply) error {
	if err := common.ValidateArgsNReply(args, reply); err != nil {
		return err
	}

	// TODO: Check later whether function needs to return anything
	s.master.handleHeartbeat(args.Address, args.FileHashes)

	return nil
}

func (s *MasterService) WriteFile(args *protocol.WriteFileRequestArgs, reply *protocol.WriteFileRequestReply) error {
	if err := common.ValidateArgsNReply(args, reply); err != nil {
		return err
	}

	workerData, err := s.master.handleWriteFileRequest(args.Filename)
	if err != nil {
		return err
	}

	reply.WorkerUrl = string(workerData)

	return nil
}
