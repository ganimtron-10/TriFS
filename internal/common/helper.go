package common

import "fmt"

func ValidateArgsNReply(args, reply any) error {
	if args == nil {
		return fmt.Errorf("rpc args is empty")
	}
	if reply == nil {
		return fmt.Errorf("rpc reply is empty")
	}

	return nil
}
