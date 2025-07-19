package worker

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

func (w *Worker) handleReadFile(filename string) ([]byte, error) {

	filenameHash := common.Hash(filename)

	w.fileStoreLock.RLock()
	fileInfo := w.fileStore[filenameHash]
	w.fileStoreLock.RUnlock()

	if fileInfo == nil {
		err := fmt.Errorf("fileinfo for file(%s) not found in worker filestore", filename)
		logger.Error(common.COMPONENT_WORKER, err.Error())
		return nil, err
	}

	// TODO: Use hashing or id gen instead of using Address
	fullFilePath := path.Join(w.Address, filename)
	file, err := os.Open(fullFilePath)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Error while opening file named %s. Error: %s", fullFilePath, err))
		return nil, err
	}
	defer file.Close()

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

func (w *Worker) handleWriteFile(filename string, data []byte) error {

	filenameHash := common.Hash(filename)
	// TODO: Add Pack Creation and Handling Logic
	w.fileStoreLock.Lock()
	w.fileStore[filenameHash] = &FileInfo{
		PackId: filenameHash,
		Offset: 0,
		Size:   len(data),
	}
	w.fileStoreLock.Unlock()

	fullFilePath := path.Join(w.Address, filename)
	if err := os.WriteFile(fullFilePath, data, 0644); err != nil {
		logger.Error(common.COMPONENT_WORKER, fmt.Sprintf("Error while writing to file named %s. Error: %s", fullFilePath, err))
		return err
	}

	return nil
}
