package master

import (
	"fmt"
	"math/rand"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
)

func (m *Master) handleGetFileWorkers(filename string) ([]string, error) {
	// return the worker url to access the file
	m.FileHashWorkerMapLock.RLock()
	defer m.FileHashWorkerMapLock.RUnlock()

	fileWorkerSet, ok := m.FileHashWorkerMap[common.GetFileHash(filename)]

	if !ok {
		return nil, fmt.Errorf("file not found")
	}

	return getList(fileWorkerSet), nil
}

func (m *Master) chooseWorkers() ([]string, error) {
	// choose a worker for writing file

	m.WorkerPoolLock.Lock()
	defer m.WorkerPoolLock.Unlock()

	workerCount := len(m.WorkerPool)
	if workerCount == 0 {
		return []string{}, fmt.Errorf("no worker available")
	}

	workers := make([]string, 0, workerCount)
	for worker := range m.WorkerPool {
		workers = append(workers, worker)
	}

	index := rand.Intn(workerCount)
	workers[0], workers[index] = workers[index], workers[0]
	return workers, nil
}

func (m *Master) handleAllocateFileWorkers(filename string) ([]string, error) {
	// choose and return the worker url to write to the file

	workerList, err := m.chooseWorkers()
	if err != nil {
		logger.Error(common.COMPONENT_MASTER, "No Worker in WorkerPool")
		return nil, err
	}

	return workerList, nil
}

func (m *Master) updateFileHashWorkerMap(workerUrl string, prevFileHashes, curFileHashes map[string]struct{}) {

	for oldFileHash := range prevFileHashes {
		if _, foundInNew := curFileHashes[oldFileHash]; !foundInNew {
			oldFileWorkerSet, ok := m.FileHashWorkerMap[oldFileHash]
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
			m.FileHashWorkerMap[newFileHash] = FileWorkerSet{}
			// TODO: Update the FileWorkerSet to properly set primary or secondary worker for now adding as secondary
			newFileWorkerSet := m.FileHashWorkerMap[newFileHash]
			newFileWorkerSet[workerUrl] = false
		}
	}
}

func (m *Master) handleHeartbeat(workerUrl string, fileHashes map[string]struct{}) {
	// update worker pool and filehash map

	m.WorkerPoolLock.Lock()
	defer m.WorkerPoolLock.Unlock()
	m.FileHashWorkerMapLock.Lock()
	defer m.FileHashWorkerMapLock.Unlock()

	workerInfo, ok := m.WorkerPool[workerUrl]
	if !ok {
		workerInfo = &WorkerInfo{}
		m.WorkerPool[workerUrl] = workerInfo
	}

	m.updateFileHashWorkerMap(workerUrl, workerInfo.FileHashes, fileHashes)
	workerInfo.FileHashes = fileHashes

}
