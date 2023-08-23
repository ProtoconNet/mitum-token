package utils

import (
	"fmt"
)

func ErrStringDecodeBSON(i interface{}) string {
	return fmt.Sprintf("failed to decode bson of %T", i)
}

func ErrStringDecodeJSON(i interface{}) string {
	return fmt.Sprintf("failed to decode json of %T", i)
}

func ErrStringUnmarshal(i interface{}) string {
	return fmt.Sprintf("failed to unmarshal %T", i)
}

func ErrStringTypeCast(expected interface{}, received interface{}) string {
	return fmt.Sprintf("expected %T, not %T", expected, received)
}
