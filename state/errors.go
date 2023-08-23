package state

import "fmt"

func ErrStringStateNotFound(k string) string {
	return fmt.Sprintf("state(key: %s) not found in state", k)
}
