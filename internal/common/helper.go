package common

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"

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
