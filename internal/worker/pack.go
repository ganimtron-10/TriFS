package worker

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ganimtron-10/TriFS/internal/common"
	"github.com/ganimtron-10/TriFS/internal/logger"
	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/klauspost/reedsolomon"
	"google.golang.org/protobuf/proto"
)

const (
	dataShards   int = 4
	parityShards int = 2
	totalShards  int = dataShards + parityShards
)

func getWalFileData(walFilePath string) ([]string, error) {
	walFileData, err := os.ReadFile(walFilePath)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to read walFile", "walFilePath", walFilePath)
		return nil, err
	}

	return strings.Split(strings.Trim(string(walFileData), "\n"), "\n"), nil
}

func createPackData(filenameHashes []string, basePath string) ([]byte, error) {
	rawDataBuffer := bytes.Buffer{}
	packInfo := protocol.PackInfo{}

	curOffset := 0
	for _, filenameHash := range filenameHashes {

		filePath := getFullFilePath(basePath, common.FOLDER_DATA, filenameHash)
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			logger.Error(common.COMPONENT_WORKER, "Unable to getFile by FilenameHash", "filenameHash", filenameHash)
			return nil, err
		}

		curSize, err := rawDataBuffer.Write(fileData)
		if err != nil {
			logger.Error(common.COMPONENT_WORKER, "Unable to write to rawDataBuffer", "error", err.Error(), "fileData", fileData)
			return nil, err
		}

		fileInfo := protocol.FileInfo{
			FilenameHash: filenameHash,
			Offset:       uint64(curOffset),
			Size:         uint64(curSize),
		}
		packInfo.FileInfos = append(packInfo.FileInfos, &fileInfo)
		curOffset += curSize
	}

	packInfoData, err := proto.Marshal(&packInfo)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to marshal packInfo", "error", err.Error(), "packInfo", packInfo)
		return nil, err
	}
	packInfoSize := uint64(len(packInfoData))

	completePackFile := []byte{}
	completePackFile = append(completePackFile, binary.AppendUvarint([]byte{}, packInfoSize)...)
	completePackFile = append(completePackFile, packInfoData...)
	completePackFile = append(completePackFile, rawDataBuffer.Bytes()...)

	return completePackFile, nil
}

func createPack(walFilePath string, basePath string) (string, []byte, error) {

	filenameHashes, err := getWalFileData(walFilePath)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to getWalLogs", "error", err.Error())
		return "", nil, err
	}

	packData, err := createPackData(filenameHashes, basePath)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to createPackData", "error", err.Error(), "filenameHashes", filenameHashes)
		return "", nil, err
	}

	packId := common.Hash(fmt.Sprintf("%s-%s", basePath, time.Now().String()))

	packFullFilePath := getFullFilePath(basePath, common.FOLDER_PACK, packId)
	err = os.WriteFile(packFullFilePath, packData, 0644)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to write to pack file", "fileName", packFullFilePath)
		return "", nil, err
	}

	return packId, packData, nil
}

func erasurePackFile(packData []byte) ([][]byte, error) {
	enc, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to create encoder", "error", err.Error())
		return nil, err
	}

	shards, err := enc.Split(packData)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to split packData", "error", err.Error(), "packData", packData)
		return nil, err
	}

	err = enc.Encode(shards)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to encode shards", "error", err.Error(), "shards", shards)
		return nil, err
	}

	return shards, nil
}

func (w *Worker) distributePackShards(shards [][]byte, packId string) error {

	// TODO: Handle mapping of these packs in worker and notify Master on Heartbeat
	shardsToKeep := 2
	for i := 0; i < shardsToKeep; i++ {
		w.handleWritePack(fmt.Sprintf("%s-shard%d", packId, i), shards[i])
	}

	ctx, cancel := context.WithTimeout(w.ctx, time.Second*5)
	defer cancel()

	conn, err := common.DialGRPC(w.MasterAddress)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Failed to connect to Master", "error", err)
		return err
	}
	defer conn.Close()

	masterClient := protocol.NewMasterServiceClient(conn)
	req := &protocol.AllocateFileWorkersRequest{Filename: packId}
	res, err := masterClient.AllocateFileWorkers(ctx, req)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Master AllocateFileWorkers error", "error", err)
		return err
	}

	if len(res.WorkerUrls) < (totalShards-shardsToKeep)/2 {
		logger.Error(common.COMPONENT_WORKER, "Not enough Workers to send Shards", "workerUrls", res.WorkerUrls)
		return fmt.Errorf("not enough workers to send shards, got only %d", len(res.WorkerUrls))
	}

	// TODO: Handle Error while sending shards
	for i, workerAddress := range res.WorkerUrls[:2] {
		go func(i int, workerAddress string) {
			ctx, cancel := context.WithTimeout(w.ctx, time.Second*5)
			defer cancel()

			conn, err := common.DialGRPC(workerAddress)
			if err != nil {
				logger.Error(common.COMPONENT_WORKER, "Failed to connect to Worker", "error", err, "workerAddress", workerAddress)
			}
			defer conn.Close()

			workerClient := protocol.NewWorkerServiceClient(conn)
			shardCount := 2 * (i + 1)
			req := &protocol.WriteRequest{Filename: fmt.Sprintf("%s-shard%d", packId, shardCount), Data: shards[shardCount]}
			res, err := workerClient.WritePack(ctx, req)
			if err != nil {
				logger.Error(common.COMPONENT_WORKER, "Worker WritePack error", "error", err)
			}
			logger.Info(common.COMPONENT_WORKER, "Successfully wrote shard", "workerAddress", workerAddress, "shard", 2*(i+1))

			req = &protocol.WriteRequest{Filename: fmt.Sprintf("%s-shard%d", packId, shardCount+1), Data: shards[(shardCount + 1)]}
			res, err = workerClient.WritePack(ctx, req)
			if err != nil {
				logger.Error(common.COMPONENT_WORKER, "Worker WritePack error", "error", err)
			}
			logger.Info(common.COMPONENT_WORKER, "Successfully wrote shard", "workerAddress", workerAddress, "shard", (2*(i+1) + 1))
			_ = res

		}(i, workerAddress)
	}

	return nil

}

func (w *Worker) startPacking(walFilePath string) error {

	packId, packData, err := createPack(walFilePath, w.Id)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to create pack file", "error", err.Error(), "walFilePath", walFilePath, "basePath", w.Id)
		return err
	}

	shards, err := erasurePackFile(packData)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to erasure pack file", "error", err.Error(), "packData", packData)
		return err
	}

	err = w.distributePackShards(shards, packId)
	if err != nil {
		logger.Error(common.COMPONENT_WORKER, "Unable to distribute pack file shards", "error", err.Error(), "shards", shards)
		return err
	}

	return nil
}
