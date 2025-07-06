package common

import (
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
