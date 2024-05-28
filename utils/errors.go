package utils

import (
	"fmt"
	"strings"
)

func ErrStringInvalid(t any) string {
	return fmt.Sprintf("invalid %T", t)
}

func ErrStringDecodeBSON(t any) string {
	return fmt.Sprintf("decode bson of %T", t)
}

func ErrStringDecodeJSON(t any) string {
	return fmt.Sprintf("decode json of %T", t)
}

func ErrStringUnPack(t any) string {
	return fmt.Sprintf("unpack %T", t)
}

func ErrStringTypeCast(expected any, received any) string {
	return fmt.Sprintf("expected %T, not %T", expected, received)
}

func ErrStringCreate(k string) string {
	return fmt.Sprintf("failed to create %s", k)
}

func ErrStringWrap(s string, e error) string {
	if e != nil {
		return fmt.Sprintf("%s: %v", s, e)
	}
	return s
}

//
//func ErrStringDuplicate(name, k string) string {
//	return fmt.Sprintf("duplicate %s found, %s", name, k)
//}

func JoinStringers(s ...fmt.Stringer) string {
	ss := make([]string, len(s))
	for i, str := range s {
		ss[i] = str.String()
	}
	return strings.Join(ss, ", ")
}
