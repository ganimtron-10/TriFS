package common

import (
	"testing"

	"github.com/ganimtron-10/TriFS/internal/protocol"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateRequest_Success(t *testing.T) {
	req := &protocol.HeartbeatRequest{}

	err := ValidateRequest(req)

	assert.Nil(t, err, "ValidateRequest should return nil when both input are not nil")
}

func TestValidateRequest_ValidateNilInput(t *testing.T) {

	err := ValidateRequest(nil)

	assert.Error(t, err, "ValidateRequest should return an error when args is nil")

	st, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, st.Code(), "Error code should be InvalidArgument for nil request")
	assert.Contains(t, st.Message(), "request is nil", "Error message should indicate nil request")

}
