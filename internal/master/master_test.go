package master

import (
	"testing"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestCreateMaster(t *testing.T) {

	master := CreateMaster()

	assert.NotNil(t, master, "CreateMaster should not return nil")
	assert.NotNil(t, master.MasterConfig, "MasterConfig should be initialized")
	assert.Equal(t, common.DEFAULT_MASTER_PORT, master.Port, "Master should use default port")
	assert.NotNil(t, master.WorkerPool, "WorkerPool should be initialized")
	assert.Empty(t, master.WorkerPool, "WorkerPool should be empty initially")

}

func TestGetDefaultMasterConfig(t *testing.T) {
	config := getDefaultMasterConfig()

	assert.NotNil(t, config, "getDefaultMasterConfig should not return nil")
	assert.Equal(t, common.DEFAULT_MASTER_PORT, config.Port, "Default config should use common.DEFAULT_MASTER_PORT")
}

func TestMaster_AddConfig(t *testing.T) {
	master := CreateMaster()
	newConfig := &MasterConfig{
		Port: 9000,
	}

	updatedMaster := master.AddConfig(newConfig)

	assert.Same(t, newConfig, updatedMaster.MasterConfig, "MasterConfig should be updated to the new config pointer")
	assert.Equal(t, 9000, updatedMaster.Port, "Port in updated config should be 9000")

	assert.Same(t, master, updatedMaster, "AddConfig should return the same master instance")
}

func TestMaster_HandleReadFile_Success(t *testing.T) {
	master := CreateMaster()

	filename := "testfile.txt"
	data, err := master.handleReadFile(filename)

	assert.NoError(t, err, "handleReadFile should not return an error on success")

	expectedData := []byte{0, 1, 2, 3, 4, 5}
	assert.Equal(t, expectedData, data, "Returned data should match expected bytes")
}

func TestMaster_HandleWriteFileRequest_NoWorkers(t *testing.T) {

	master := CreateMaster()

	filename := "newfile.txt"
	workerURL, err := master.handleWriteFileRequest(filename)

	assert.Error(t, err, "handleWriteFileRequest should return an error when no workers are available")
	assert.EqualError(t, err, "worker not available. please try later", "Error message should match expected")

	assert.Nil(t, workerURL, "WorkerURL should be nil when no workers are available")

}

func TestMaster_HandleWriteFileRequest_OneWorker(t *testing.T) {

	master := CreateMaster()

	master.WorkerPoolLock.Lock()
	master.WorkerPool["worker-1:9000"] = 1
	master.WorkerPoolLock.Unlock()

	filename := "anotherfile.txt"
	workerURL, err := master.handleWriteFileRequest(filename)

	assert.NoError(t, err, "handleWriteFileRequest should not return an error when workers are available")

	expectedWorkerURL := []byte("worker-1:9000")
	assert.Equal(t, expectedWorkerURL, workerURL, "Returned worker URL should match the single worker in the pool")

}

func TestMaster_HandleWriteFileRequest_MultipleWorkers(t *testing.T) {

	master := CreateMaster()

	master.WorkerPoolLock.Lock()
	master.WorkerPool["worker-A:9000"] = 1
	master.WorkerPool["worker-B:9001"] = 2
	master.WorkerPool["worker-C:9002"] = 3
	master.WorkerPoolLock.Unlock()

	filename := "multi_worker_file.txt"
	workerURL, err := master.handleWriteFileRequest(filename)

	assert.NoError(t, err, "handleWriteFileRequest should not return an error when multiple workers are available")

	expectedWorkers := []string{"worker-A:9000", "worker-B:9001", "worker-C:9002"}
	returnedWorkerStr := string(workerURL)
	assert.Contains(t, expectedWorkers, returnedWorkerStr, "Returned worker URL should be one of the added workers")

}
