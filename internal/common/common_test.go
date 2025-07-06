package common

import (
	"testing"

	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/stretchr/testify/assert"
)

func TestValidateArgsNReply_Success(t *testing.T) {
	args := &protocol.HeartBeatArgs{}
	reply := &protocol.HeartBeatReply{}

	err := ValidateArgsNReply(args, reply)

	assert.Nil(t, err, "ValidateArgsNReply should return nil when both input are not nil")
}

func TestValidateArgsNReply_ValidateNilInput(t *testing.T) {
	args := &protocol.HeartBeatArgs{}
	reply := &protocol.HeartBeatReply{}

	err := ValidateArgsNReply(nil, reply)

	assert.Error(t, err, "ValidateArgsNReply should return an error when args is nil")
	assert.EqualError(t, err, "rpc args is empty", "Error message should indicate empty args")

	err = ValidateArgsNReply(args, nil)

	assert.Error(t, err, "ValidateArgsNReply should return an error when reply is nil")
	assert.EqualError(t, err, "rpc reply is empty", "Error message should indicate empty reply")
}

func TestValidateArgsNReply_ValidateNilPointer(t *testing.T) {
	var argsPtr *protocol.HeartBeatArgs
	var replyPtr *protocol.HeartBeatReply

	err := ValidateArgsNReply(argsPtr, replyPtr)

	assert.Error(t, err, "ValidateArgsNReply should return an error when args points to nil")
	assert.EqualError(t, err, "rpc args is empty", "Error message should indicate empty args")

	args := &protocol.HeartBeatArgs{}
	err = ValidateArgsNReply(args, replyPtr)

	assert.Error(t, err, "ValidateArgsNReply should return an error when reply points to nil")
	assert.EqualError(t, err, "rpc reply is empty", "Error message should indicate empty reply")
}
