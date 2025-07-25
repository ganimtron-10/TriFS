package worker

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

type WAL struct {
	Logs     []string
	BasePath string
	WALLock  sync.RWMutex
}

func createWAL(basePath string) WAL {
	return WAL{
		Logs:     []string{},
		BasePath: basePath,
	}
}

func (wal *WAL) getWALFilePath() string {
	fileName := fmt.Sprintf("wal-%s.log", common.Hash(time.Now().String()))
	return getFullFilePath(wal.BasePath, common.FOLDER_WAL, fileName)
}

func (wal *WAL) Clear() {
	wal.Logs = []string{}
}

func (wal *WAL) addLog(filenameHash string) {
	wal.WALLock.Lock()
	wal.Logs = append(wal.Logs, filenameHash)
	wal.WALLock.Unlock()

	// Simulating Flushing
	if len(wal.Logs) > 2 {
		copiedLogs, err := wal.flushToFile()
		if err != nil {
			logger.Error(common.COMPONENT_WORKER, "Unable to flush WAL", "error", err.Error())
			// TODO: Need to handle this error
		}
		wal.Clear()

		// Start Packing
		_ = copiedLogs
	}
}

func (wal *WAL) flushToFile() ([]string, error) {
	if len(wal.Logs) == 0 {
		logger.Info(common.COMPONENT_WORKER, "Nothing to Flush, WAL is Empty")
		return nil, nil
	}

	copiedLogs := make([]string, len(wal.Logs))

	wal.WALLock.RLock()
	copy(copiedLogs, wal.Logs)
	wal.WALLock.RUnlock()

	walFilePath := wal.getWALFilePath()
	logger.Info(common.COMPONENT_WORKER, "Retrieved FilePath", "path", walFilePath)
	if err := os.MkdirAll(filepath.Dir(walFilePath), 0644); err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to create WAL directory", "error", err.Error(), "path", walFilePath)
		return nil, err
	}

	stringifiedLogs := ""
	for _, log := range copiedLogs {
		stringifiedLogs += log + "\n"
	}

	err := os.WriteFile(walFilePath, []byte(stringifiedLogs), 0644)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to write to wal file", "error", err.Error(), "path", walFilePath, "data", stringifiedLogs)
		return nil, err
	}

	return copiedLogs, nil
}
