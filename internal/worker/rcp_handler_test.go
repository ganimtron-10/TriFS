package worker

import (
	"os"
	"path"
	"testing"

	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/stretchr/testify/assert"
)

func createTestWorkerAndService() (*Worker, *WorkerService) {
	worker, _ := CreateWorker()
	service := CreateWorkerService(worker)
	return worker, service
}

func TestCreateWorkerService(t *testing.T) {
	dummyWorker := &Worker{}

	service := CreateWorkerService(dummyWorker)

	assert.NotNil(t, service, "CreateWorkerService should not return nil")
	assert.Same(t, dummyWorker, service.worker, "WorkerService should hold a reference to the provided Worker instance")
}

func TestWorkerService_WriteFile_Success(t *testing.T) {
	worker, service := createTestWorkerAndService()

	filename := "testfile.txt"
	data := []byte("This is some test data for the file.")

	args := &protocol.WriteFileArgs{Filename: filename, Data: data}
	reply := &protocol.WriteFileReply{}

	err := service.WriteFile(args, reply)

	assert.NoError(t, err, "WriteFile should not return an error on success")

	fullFilePath := path.Join(worker.Address, filename)
	assert.FileExists(t, fullFilePath, "File should be created on disk")

	readData, err := os.ReadFile(fullFilePath)
	assert.NoError(t, err, "Should be able to read the file from disk")
	assert.Equal(t, data, readData, "Content of the file on disk should match original data")

	filenameHash := hash(filename)
	fileInfoEntry, exists := worker.fileStore[filenameHash]
	assert.True(t, exists, "File info should be added to fileStore")
	assert.NotNil(t, fileInfoEntry, "FileInfo entry should not be nil in fileStore")
	assert.Equal(t, filenameHash, fileInfoEntry.PackId, "PackId in FileInfo should match filename hash")
	assert.Equal(t, 0, fileInfoEntry.Offset, "Offset in FileInfo should be 0 (as hardcoded in handleWriteFile)")

	// Cleanup
	os.RemoveAll(worker.Address)
}
