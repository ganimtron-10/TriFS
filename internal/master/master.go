package master

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

type MasterConfig struct {
	Port int
}

type WorkerInfo struct {
	FileHashes map[string]struct{}
}

type FileWorkerSet = map[string]bool

func getList(fileWorkerSet FileWorkerSet) []string {
	workerList := make([]string, 0, len(fileWorkerSet))
	for worker := range fileWorkerSet {
		workerList = append(workerList, worker)
	}
	return workerList
}

type Master struct {
	*MasterConfig
	WorkerPool            map[string]*WorkerInfo
	WorkerPoolLock        sync.RWMutex
	FileHashWorkerMap     map[string]FileWorkerSet
	FileHashWorkerMapLock sync.RWMutex
}

func getDefaultMasterConfig() *MasterConfig {
	return &MasterConfig{
		Port: common.DEFAULT_MASTER_PORT,
	}
}

func CreateMaster() *Master {
	logger.Info(common.COMPONENT_MASTER, "Creating Master...")
	return &Master{
		MasterConfig:      getDefaultMasterConfig(),
		WorkerPool:        make(map[string]*WorkerInfo),
		FileHashWorkerMap: make(map[string]FileWorkerSet),
	}
}

func (master *Master) AddConfig(config *MasterConfig) *Master {
	master.MasterConfig = config
	return master
}

func (master *Master) handleReadFile(filename string) ([]string, error) {
	// return the worker url to access the file
	master.FileHashWorkerMapLock.RLock()
	defer master.FileHashWorkerMapLock.RUnlock()

	fileWorkerSet, ok := master.FileHashWorkerMap[common.Hash(filename)]

	if !ok {
		return nil, fmt.Errorf("file not found")
	}

	return getList(fileWorkerSet), nil
}

func (master *Master) chooseWorker() (string, error) {
	// choose a worker for writing file

	workerCount := len(master.WorkerPool)
	if workerCount == 0 {
		return "", fmt.Errorf("no worker available")
	}

	workers := make([]string, 0, workerCount)
	for worker := range master.WorkerPool {
		workers = append(workers, worker)
	}
	return workers[rand.Intn(workerCount)], nil
}

func (master *Master) handleWriteFileRequest(filename string) ([]byte, error) {
	// choose and return the worker url to write to the file

	master.WorkerPoolLock.Lock()
	defer master.WorkerPoolLock.Unlock()

	worker, err := master.chooseWorker()
	if err != nil {
		logger.Error(common.COMPONENT_MASTER, "No Worker in WorkerPool")
		return nil, err
	}

	return []byte(worker), nil
}

func (master *Master) updateFileHashWorkerMap(workerUrl string, prevFileHashes, curFileHashes map[string]struct{}) {

	for oldFileHash := range prevFileHashes {
		if _, foundInNew := curFileHashes[oldFileHash]; !foundInNew {
			oldFileWorkerSet, ok := master.FileHashWorkerMap[oldFileHash]
			if !ok {
				logger.Error(common.COMPONENT_MASTER, "OldFileHash entry not present in FileHashWorkerMap", "OldFileHash", oldFileHash, "WorkerUrl", workerUrl)
				// TODO: This shouldnt happen, Either Error out or add a way to handle the condition
			}

			isPrimary, ok := oldFileWorkerSet[workerUrl]
			if !ok {
				logger.Error(common.COMPONENT_MASTER, "WorkerUrl entry not present in FileWorkerSet", "WorkerUrl", workerUrl, "FileWorkerSet", oldFileWorkerSet)
				// TODO: This shouldnt happen, plan a way to handle this condition
			}

			// TODO: This file doesnt have a primary worker, choose a primary worker
			if isPrimary {
			}
			delete(oldFileWorkerSet, workerUrl)
			// TODO: Trigger a replication
		}
	}

	for newFileHash := range curFileHashes {
		if _, foundInOld := prevFileHashes[newFileHash]; !foundInOld {
			master.FileHashWorkerMap[newFileHash] = FileWorkerSet{}
			// TODO: Update the FileWorkerSet to properly set primary or secondary worker for now adding as secondary
			newFileWorkerSet := master.FileHashWorkerMap[newFileHash]
			newFileWorkerSet[workerUrl] = false
		}
	}
}

func (master *Master) handleHeartbeat(workerUrl string, fileHashes map[string]struct{}) {
	// update worker pool and filehash map

	master.WorkerPoolLock.Lock()
	defer master.WorkerPoolLock.Unlock()
	master.FileHashWorkerMapLock.Lock()
	defer master.FileHashWorkerMapLock.Unlock()

	workerInfo, ok := master.WorkerPool[workerUrl]
	if !ok {
		workerInfo = &WorkerInfo{}
		master.WorkerPool[workerUrl] = workerInfo
	}

	master.updateFileHashWorkerMap(workerUrl, workerInfo.FileHashes, fileHashes)
	workerInfo.FileHashes = fileHashes

}
