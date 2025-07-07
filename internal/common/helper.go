package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
)

func ValidateArgsNReply(args, reply any) error {
	if args == nil || reflect.ValueOf(args).IsNil() {
		return fmt.Errorf("rpc args is empty")
	}
	if reply == nil || reflect.ValueOf(reply).IsNil() {
		return fmt.Errorf("rpc reply is empty")
	}

	return nil
}

func Hash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
