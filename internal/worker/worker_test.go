package worker

import (
	"context"
	"os"
	"testing"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetDefaultWorkerConfig(t *testing.T) {
	config := getDefaultWorkerConfig()

	assert.NotNil(t, config, "getDefaultWorkerConfig should not return nil")
	assert.Equal(t, common.DEFAULT_MASTER_ADDRESS, config.MasterAddress, "Default MasterAddress should be common.DEFAULT_MASTER_ADDRESS")
	assert.NotEmpty(t, config.Address, "Default Address should not be empty")
	assert.Equal(t, 5, config.HeartbeatInterval, "Default HeartbeatInterval should be 5")
}

func TestCreateWorker(t *testing.T) {
	worker, err := createWorker()
	defer os.RemoveAll(worker.Address)

	assert.NoError(t, err, "createWorker should not return an error")
	assert.NotNil(t, worker, "createWorker should not return nil Worker")
	assert.NotNil(t, worker.fileStore, "fileStore should be initialized")
	assert.DirExists(t, worker.Address, "Worker directory should be created")

}

func TestWorker_AddConfig(t *testing.T) {
	worker, _ := createWorker()
	defer os.RemoveAll(worker.Address)

	newConfig := &WorkerConfig{
		MasterAddress:     "new-master:9000",
		Address:           "new-worker:9001",
		HeartbeatInterval: 10,
	}

	updatedWorker := worker.AddConfig(newConfig)

	assert.Same(t, newConfig, updatedWorker.WorkerConfig, "WorkerConfig should be updated to the new config pointer")
	assert.Equal(t, "new-master:9000", updatedWorker.MasterAddress, "MasterAddress should be updated to new-master:9000")
	assert.Equal(t, "new-worker:9001", updatedWorker.Address, "Address should be updated to new-worker:9001")
	assert.Equal(t, 10, updatedWorker.HeartbeatInterval, "HeartbeatInterval should be updated to 10")

	assert.Same(t, worker, updatedWorker, "AddConfig should return the same worker instance")

	// Cleanup directory created during test
	os.RemoveAll(worker.Address)
}

func TestGetFileHashes(t *testing.T) {
	fileStore := map[string]*FileInfo{
		"hash1": {PackId: "p1", Offset: 0, Size: 10},
		"hash2": {PackId: "p2", Offset: 10, Size: 20},
	}

	hashes := getFileHashes(fileStore)

	assert.Len(t, hashes, 2, "There should be exactly 2 file hashes returned")
	assert.Contains(t, hashes, "hash1", "File hashes should include 'hash1'")
	assert.Contains(t, hashes, "hash2", "File hashes should include 'hash2'")
}

func TestReadFile_ValidationError(t *testing.T) {
	worker, _ := createWorker()
	defer os.RemoveAll(worker.Address)

	_, err := worker.ReadFile(context.Background(), nil)
	assert.Error(t, err, "ReadFile should return error when request is nil")
	st, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, st.Code(), "Error code should be InvalidArgument for nil request")
}

func TestWriteFile_ValidationError(t *testing.T) {
	worker, _ := createWorker()
	defer os.RemoveAll(worker.Address)

	_, err := worker.WriteFile(context.Background(), nil)
	assert.Error(t, err, "WriteFile should return error when request is nil")
	st, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, st.Code(), "Error code should be InvalidArgument for nil request")
}

func TestReadFile_Success(t *testing.T) {
	worker, _ := createWorker()
	defer os.RemoveAll(worker.Address)

	expectedFilename := "exists.txt"
	expectedFileData := []byte("File content inside the exists.txt")
	worker.handleWriteFile(expectedFilename, expectedFileData)

	req := &protocol.ReadFileRequest{Filename: expectedFilename}

	res, err := worker.ReadFile(context.Background(), req)

	assert.NoError(t, err, "ReadFile should not return error when file exists")
	assert.NotNil(t, res, "Response should not be nil when file exists")
	assert.Equal(t, expectedFilename, res.Filename, "Response Filename should match request Filename")
	assert.Equal(t, expectedFileData, res.Data, "Response Data should match file content")
}

func TestReadFile_FileNotFound(t *testing.T) {
	worker, _ := createWorker()
	defer os.RemoveAll(worker.Address)

	req := &protocol.ReadFileRequest{Filename: "missing.txt"}
	res, err := worker.ReadFile(context.Background(), req)

	assert.Error(t, err, "ReadFile should return an error when file is missing")
	assert.Nil(t, res, "Response should be nil when file is not found")

	st, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.Internal, st.Code(), "Error code should be Internal when file is missing")
	assert.Contains(t, st.Message(), "unable to read file", "Error message should indicate inability to read file")
}

func TestWriteFile_Success(t *testing.T) {
	worker, _ := createWorker()
	defer os.RemoveAll(worker.Address)

	req := &protocol.WriteFileRequest{Filename: "writable.txt", Data: []byte("some data")}
	res, err := worker.WriteFile(context.Background(), req)

	assert.NoError(t, err, "WriteFile should not return error on success")
	assert.NotNil(t, res, "Response should not be nil on successful write")
}

func TestWriteFile_Error(t *testing.T) {
	worker, _ := createWorker()
	defer os.RemoveAll(worker.Address)

	// simulating error: removing base directory
	os.RemoveAll(worker.Address)

	req := &protocol.WriteFileRequest{Filename: "readonly.txt", Data: []byte("some data")}
	res, err := worker.WriteFile(context.Background(), req)

	assert.Error(t, err, "WriteFile should return an error when write fails")
	assert.Nil(t, res, "Response should be nil when write fails")

	st, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.Internal, st.Code(), "Error code should be Internal for write failure")
	assert.Contains(t, st.Message(), "unable to write file", "Error message should indicate inability to write file")
}
