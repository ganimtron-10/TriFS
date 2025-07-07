package worker

import (
	"fmt"
	"io"
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

type FileInfo struct {
	PackId string
	Offset int
	Size   int
}

type Worker struct {
	*WorkerConfig
	fileStore map[string]*FileInfo
}

func getDefaultWorkerConfig() *WorkerConfig {
	return &WorkerConfig{
		MasterAddress:     common.DEFAULT_MASTER_ADDRESS,
		Address:           transport.GetAddressWithRandomPort(),
		HeartbeatInterval: 5,
	}
}

func getFileHashes(fileStore map[string]*FileInfo) map[string]struct{} {
	fileHashes := make(map[string]struct{})
	for fileHash := range fileStore {
		fileHashes[fileHash] = struct{}{}
	}
	return fileHashes
}

func (w *Worker) SendHeartBeat() {
	transport.DialRpcCall(w.MasterAddress, "MasterService.HeartBeat", &protocol.HeartBeatArgs{Address: w.Address, FileHashes: getFileHashes(w.fileStore)}, &protocol.HeartBeatReply{})
}

func CreateWorker() (*Worker, error) {
	logger.Info(common.COMPONENT_WORKER, "Creating Worker...")

	worker := &Worker{
		WorkerConfig: getDefaultWorkerConfig(),
		fileStore:    make(map[string]*FileInfo),
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

func (worker *Worker) handleReadFile(filename string) ([]byte, error) {

	filenameHash := common.Hash(filename)
	fileInfo := worker.fileStore[filenameHash]

	fullFilePath := path.Join(worker.Address, filename)
	file, err := os.Open(fullFilePath)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Error while opening file named %s. Error: %s", fullFilePath, err))
		return nil, err
	}

	_, err = file.Seek(int64(fileInfo.Offset), io.SeekStart)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Error while seeking file named %s. Error: %s", fullFilePath, err))
		return nil, err
	}

	fileData := make([]byte, fileInfo.Size)

	_, err = io.ReadFull(file, fileData)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Error while reading from file named %s. Error: %s", fullFilePath, err))
		return nil, err
	}

	return fileData, nil
}

func (worker *Worker) handleWriteFile(filename string, data []byte) error {

	filenameHash := common.Hash(filename)
	// TODO: Add Pack Creation and Handling Logic
	worker.fileStore[filenameHash] = &FileInfo{
		PackId: filenameHash,
		Offset: 0,
		Size:   len(data),
	}

	fullFilePath := path.Join(worker.Address, filename)
	if err := os.WriteFile(fullFilePath, data, 0755); err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Error while writing to file named %s. Error: %s", fullFilePath, err))
		return err
	}

	return nil
}
