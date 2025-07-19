package master

import (
	"testing"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateMaster(t *testing.T) {
	master := createMaster()

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
	master := createMaster()
	newConfig := &MasterConfig{
		Port: 9000,
	}

	updatedMaster := master.AddConfig(newConfig)

	assert.Same(t, newConfig, updatedMaster.MasterConfig, "MasterConfig should be updated to the new config pointer")
	assert.Equal(t, 9000, updatedMaster.Port, "Port in updated config should be 9000")

	assert.Same(t, master, updatedMaster, "AddConfig should return the same master instance")
}

func TestMasterService_GetFileWorkers_Success(t *testing.T) {
	master := createMaster()

	filename := "testfile_read.txt"

	master.FileHashWorkerMapLock.Lock()
	master.FileHashWorkerMap[common.Hash(filename)] = FileWorkerSet{"worker-A:9000": true}
	master.FileHashWorkerMapLock.Unlock()

	req := &protocol.GetFileWorkersRequest{Filename: filename}

	res, err := master.GetFileWorkers(t.Context(), req)

	assert.NoError(t, err, "GetFileWorkers should not return an error on success")
	assert.NotNil(t, res, "Response should not be nil on success")

	expectedData := []string{"worker-A:9000"}
	assert.Equal(t, expectedData, res.WorkerUrls, "Reply data should match expected data from handleReadFile")
}

func TestMasterService_GetFileWorkers_FileNotFound(t *testing.T) {
	master := createMaster()

	filename := "nonexistent_file.txt"
	req := &protocol.GetFileWorkersRequest{Filename: filename}

	res, err := master.GetFileWorkers(t.Context(), req)

	assert.Error(t, err, "GetFileWorkers should return an error when file is not found")
	assert.Nil(t, res, "Response should be nil on error")

	st, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.Internal, st.Code(), "Error code should be Internal for file not found")
	assert.Contains(t, st.Message(), "file not found", "Error message should indicate file not found")
}

func TestMasterService_GetFileWorkers_ValidationNilRequest(t *testing.T) {
	master := createMaster()

	res, err := master.GetFileWorkers(t.Context(), nil)

	assert.Error(t, err, "GetFileWorkers should return an error when request is nil")
	assert.Nil(t, res, "Response should be nil when request is nil")

	st, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, st.Code(), "Error code should be InvalidArgument for nil request")
	assert.Contains(t, st.Message(), "request is nil", "Error message should indicate nil request")
}

func TestMasterService_Heartbeat_Success(t *testing.T) {
	master := createMaster()

	workerAddress := "worker-heartbeat-1:9000"

	req := &protocol.HeartbeatRequest{
		WorkerAddress:    workerAddress,
		HostedFileHashes: []string{"hash1", "hash2"},
	}

	assert.Empty(t, master.WorkerPool, "WorkerPool should be empty before heartbeat")

	res, err := master.Heartbeat(t.Context(), req)

	assert.NoError(t, err, "Heartbeat should not return an error on success")
	assert.NotNil(t, res, "Response should not be nil on success")

	master.WorkerPoolLock.RLock()
	workerInfo, exists := master.WorkerPool[workerAddress]
	master.WorkerPoolLock.RUnlock()
	assert.True(t, exists, "Worker should be added to WorkerPool after heartbeat")
	assert.Equal(t, 1, len(master.WorkerPool), "WorkerPool should contain exactly one worker")

	assert.NotNil(t, workerInfo.FileHashes, "WorkerInfo FileHashes should not be nil")
	assert.Contains(t, workerInfo.FileHashes, "hash1", "WorkerInfo should contain hosted file hash")
	assert.Contains(t, workerInfo.FileHashes, "hash2", "WorkerInfo should contain hosted file hash")
	assert.Equal(t, 2, len(workerInfo.FileHashes), "WorkerInfo should have correct number of hosted file hashes")
}

func TestMasterService_Heartbeat_ValidationNilRequest(t *testing.T) {
	master := createMaster()

	res, err := master.Heartbeat(t.Context(), nil)

	assert.Error(t, err, "Heartbeat should return an error when request is nil")
	assert.Nil(t, res, "Response should be nil when request is nil")

	st, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, st.Code(), "Error code should be InvalidArgument for nil request")
	assert.Contains(t, st.Message(), "request is nil", "Error message should indicate nil request")
}

func TestMasterService_AllocateWriteLocations_Success(t *testing.T) {
	master := createMaster()

	filename := "testfile_write.txt"
	expectedWorkerList := []string{"worker-A:9000"}

	master.WorkerPoolLock.Lock()
	master.WorkerPool[expectedWorkerList[0]] = &WorkerInfo{}
	master.WorkerPoolLock.Unlock()

	req := &protocol.AllocateFileWorkersRequest{Filename: filename}

	res, err := master.AllocateFileWorkers(t.Context(), req)

	assert.NoError(t, err, "AllocateWriteLocations should not return an error on success")
	assert.NotNil(t, res, "Response should not be nil on success")
	assert.Equal(t, expectedWorkerList, res.WorkerUrls, "Reply WorkerUrls should match expected from handleWriteFileRequest")
}

func TestMasterService_AllocateWriteLocations_ValidationNilRequest(t *testing.T) {
	master := createMaster()

	res, err := master.AllocateFileWorkers(t.Context(), nil)

	assert.Error(t, err, "AllocateWriteLocations should return an error when request is nil")
	assert.Nil(t, res, "Response should be nil when request is nil")

	st, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, st.Code(), "Error code should be InvalidArgument for nil request")
	assert.Contains(t, st.Message(), "request is nil", "Error message should indicate nil request")
}

func TestMasterService_AllocateWriteLocations_NoWorkersError(t *testing.T) {
	master := createMaster()

	master.WorkerPoolLock.Lock()
	master.WorkerPool = make(map[string]*WorkerInfo)
	master.WorkerPoolLock.Unlock()

	filename := "testfile_write_no_worker.txt"
	req := &protocol.AllocateFileWorkersRequest{Filename: filename}

	res, err := master.AllocateFileWorkers(t.Context(), req)

	assert.Error(t, err, "AllocateWriteLocations should return an error when handleWriteFileRequest fails due to no workers")
	assert.Nil(t, res, "Response should be nil on error")

	st, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.Internal, st.Code(), "Error code should be Internal for no worker available")
	assert.Contains(t, st.Message(), "no worker available", "Error message should match expected from handleWriteFileRequest")
}
