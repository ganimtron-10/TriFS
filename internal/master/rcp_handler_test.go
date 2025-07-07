package master

import (
	"testing"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/stretchr/testify/assert"
)

func createTestMasterAndService() (*Master, *MasterService) {
	master := CreateMaster()
	service := CreateMasterService(master)
	return master, service
}

func TestCreateMasterService(t *testing.T) {
	dummyMaster := &Master{}

	service := CreateMasterService(dummyMaster)

	assert.NotNil(t, service, "CreateMasterService should not return nil")
	assert.Same(t, dummyMaster, service.master, "MasterService should hold a reference to the provided Master instance")
}

func TestMasterService_ReadFile_Success(t *testing.T) {
	master, service := createTestMasterAndService()

	filename := "testfile_read.txt"

	master.FileHashWorkerMapLock.Lock()
	master.FileHashWorkerMap[common.Hash(filename)] = FileWorkerSet{"worker-A:9000": true}
	master.FileHashWorkerMapLock.Unlock()

	args := &protocol.ReadFileRequestArgs{Filename: filename}
	reply := &protocol.ReadFileRequestReply{}

	err := service.ReadFile(args, reply)

	assert.NoError(t, err, "ReadFile should not return an error on success")

	expectedData := []string{"worker-A:9000"}
	assert.Equal(t, expectedData, reply.WorkerUrls, "Reply data should match expected data from handleReadFile")
}

func TestMasterService_ReadFile_Failure(t *testing.T) {
	_, service := createTestMasterAndService()

	filename := "testfile_read.txt"
	args := &protocol.ReadFileRequestArgs{Filename: filename}
	reply := &protocol.ReadFileRequestReply{}

	err := service.ReadFile(args, reply)

	assert.Error(t, err, "ReadFile should return an error on failure")
	assert.EqualError(t, err, "file not found", "Error message should match expected")
}

func TestMasterService_ReadFile_ValidationNilArgs(t *testing.T) {
	_, service := createTestMasterAndService()

	reply := &protocol.ReadFileRequestReply{}
	err := service.ReadFile(nil, reply)

	assert.Error(t, err, "ReadFile should return an error when args is nil")
	assert.EqualError(t, err, "rpc args is empty", "Error message should indicate empty arguments")
}

func TestMasterService_ReadFile_ValidationNilReply(t *testing.T) {
	_, service := createTestMasterAndService()

	args := &protocol.ReadFileRequestArgs{Filename: "test.txt"}
	err := service.ReadFile(args, nil)

	assert.Error(t, err, "ReadFile should return an error when reply is nil")
	assert.EqualError(t, err, "rpc reply is empty", "Error message should indicate empty reply")
}

func TestMasterService_HeartBeat_Success(t *testing.T) {
	master, service := createTestMasterAndService()

	workerAddress := "worker-heartbeat-1:9000"
	args := &protocol.HeartBeatArgs{Address: workerAddress}
	reply := &protocol.HeartBeatReply{}

	assert.Empty(t, master.WorkerPool, "WorkerPool should be empty before heartbeat")

	err := service.HeartBeat(args, reply)

	assert.NoError(t, err, "HeartBeat should not return an error on success")

	master.WorkerPoolLock.RLock()
	_, exists := master.WorkerPool[workerAddress]
	master.WorkerPoolLock.RUnlock()
	assert.True(t, exists, "Worker should be added to WorkerPool after heartbeat")
	assert.Equal(t, 1, len(master.WorkerPool), "WorkerPool should contain exactly one worker")
}

func TestMasterService_HeartBeat_ValidationNilArgs(t *testing.T) {
	_, service := createTestMasterAndService()

	reply := &protocol.HeartBeatReply{}
	err := service.HeartBeat(nil, reply)

	assert.Error(t, err, "HeartBeat should return an error when args is nil")
	assert.EqualError(t, err, "rpc args is empty", "Error message should indicate empty arguments")
}

func TestMasterService_HeartBeat_ValidationNilReply(t *testing.T) {
	_, service := createTestMasterAndService()

	args := &protocol.HeartBeatArgs{Address: "worker-test:9000"}
	err := service.HeartBeat(args, nil)

	assert.Error(t, err, "HeartBeat should return an error when reply is nil")
	assert.EqualError(t, err, "rpc reply is empty", "Error message should indicate empty reply")
}

func TestMasterService_WriteFile_Success(t *testing.T) {
	master, service := createTestMasterAndService()

	filename := "testfile_write.txt"
	expectedWorkerList := []string{"worker-A:9000"}

	master.WorkerPoolLock.Lock()
	master.WorkerPool[expectedWorkerList[0]] = &WorkerInfo{}
	master.WorkerPoolLock.Unlock()

	args := &protocol.WriteFileRequestArgs{Filename: filename}
	reply := &protocol.WriteFileRequestReply{}

	err := service.WriteFile(args, reply)

	assert.NoError(t, err, "WriteFile should not return an error on success")
	assert.Equal(t, expectedWorkerList, reply.WorkerUrls, "Reply WorkerUrl should match expected from handleWriteFileRequest")
}

func TestMasterService_WriteFile_ValidationNilArgs(t *testing.T) {
	_, service := createTestMasterAndService()

	reply := &protocol.WriteFileRequestReply{}
	err := service.WriteFile(nil, reply)

	assert.Error(t, err, "WriteFile should return an error when args is nil")
	assert.EqualError(t, err, "rpc args is empty", "Error message should indicate empty arguments")
}

func TestMasterService_WriteFile_ValidationNilReply(t *testing.T) {
	_, service := createTestMasterAndService()

	args := &protocol.WriteFileRequestArgs{Filename: "test.txt"}
	err := service.WriteFile(args, nil)

	assert.Error(t, err, "WriteFile should return an error when reply is nil")
	assert.EqualError(t, err, "rpc reply is empty", "Error message should indicate empty reply")
}

func TestMasterService_WriteFile_NoWorkersError(t *testing.T) {
	master, service := createTestMasterAndService()

	master.WorkerPoolLock.Lock()
	master.WorkerPool = make(map[string]*WorkerInfo)
	master.WorkerPoolLock.Unlock()

	filename := "testfile_write_no_worker.txt"
	args := &protocol.WriteFileRequestArgs{Filename: filename}
	reply := &protocol.WriteFileRequestReply{}

	err := service.WriteFile(args, reply)

	assert.Error(t, err, "WriteFile should return an error when handleWriteFileRequest fails due to no workers")
	assert.EqualError(t, err, "no worker available", "Error message should match expected from handleWriteFileRequest")
	assert.Empty(t, reply.WorkerUrls, "Reply WorkerUrl should be empty on error")
}
