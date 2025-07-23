package worker

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"google.golang.org/protobuf/proto"
)

type WAL struct {
	protocol.WAL

	BasePath string
	WALLock  sync.RWMutex
}

func createWAL(basePath string) WAL {
	return WAL{
		WAL:      protocol.WAL{},
		BasePath: basePath,
	}
}

func (wal *WAL) getWALFilePath() string {
	return path.Join(wal.BasePath, "wal", fmt.Sprintf("wal-%s.log", time.Now().Format("15-04-05")))
}

func (wal *WAL) Clear() {
	wal.WAL = protocol.WAL{}
}

func (wal *WAL) addLog(filename string, data []byte) {
	wal.WALLock.Lock()
	wal.WAL.Logs = append(wal.WAL.Logs, &protocol.FileLog{
		Filename: filename,
		Data:     data,
	})
	wal.WALLock.Unlock()

	// Simulating Flushing
	if len(wal.WAL.Logs) > 2 {
		err := wal.flushToFile()
		if err != nil {
			logger.Error(common.COMPONENT_WORKER, "Unable to flush WAL", "error", err.Error())
			// TODO: Need to handle this error
		}

	}
}

func (wal *WAL) flushToFile() error {
	wal.WALLock.RLock()
	data, err := proto.Marshal(&wal.WAL)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to marshal WAL", "error", err.Error())
		return err
	}
	wal.WALLock.RUnlock()

	filePath := wal.getWALFilePath()
	logger.Info(common.COMPONENT_WORKER, "Retrieved FilePath", "path", filePath)
	if err := os.MkdirAll(filepath.Dir(filePath), 0644); err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to create WAL directory", "error", err.Error(), "path", filePath)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to write to wal file", "error", err.Error(), "path", filePath, "data", data)
	}

	return nil
}
