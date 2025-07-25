package common

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"reflect"

	"github.com/ganimtron-10/TriFS/internal/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateRequest(req any) error {
	if req == nil || reflect.ValueOf(req).IsNil() {
		return status.Errorf(codes.InvalidArgument, "request is nil")
	}

	return nil
}

func Hash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetAddressWithRandomPort() string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		logger.Error(COMPONENT_COMMON, "Unable to creating listener", "error", err.Error())
		return ""
	}
	defer listener.Close()

	return listener.Addr().String()
}
