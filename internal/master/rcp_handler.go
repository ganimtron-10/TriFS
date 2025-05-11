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

func (s *MasterService) ReadFile(args *protocol.ReadFileArgs, reply *protocol.ReadFileReply) error {

	if args == nil {
		return fmt.Errorf("RPC Args is empty")
	}
	if reply == nil {
		return fmt.Errorf("RPC Reply is empty")
	}

	fileData, err := s.master.handleReadFile(args.Filename)
	if err != nil {
		return err
	}

	reply.Filename = args.Filename
	reply.Data = fileData

	return nil
}
