package common

import "fmt"

func ValidateArgsNReply(args, reply any) error {
	if args == nil {
		return fmt.Errorf("RPC Args is empty")
	}
	if reply == nil {
		return fmt.Errorf("RPC Reply is empty")
	}

	return nil
}
