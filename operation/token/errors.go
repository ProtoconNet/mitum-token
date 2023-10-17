package token

import (
	"fmt"

	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
)

func ErrStringPreProcess(t any) string {
	return fmt.Sprintf("failed to preprocess %T", t)
}

func ErrStringProcess(t any) string {
	return fmt.Sprintf("failed to process %T", t)
}

func ErrBaseOperationProcess(e error, formatter string, args ...interface{}) base.BaseOperationProcessReasonError {
	return base.NewBaseOperationProcessReasonError(utils.ErrStringWrap(fmt.Sprintf(formatter, args...), e))
}

func ErrStateNotFound(name string, k string, e error) base.BaseOperationProcessReasonError {
	return base.NewBaseOperationProcessReasonError(utils.ErrStringWrap(fmt.Sprintf("%s not found, %s", name, k), e))
}

func ErrStateAlreadyExists(name, k string, e error) base.BaseOperationProcessReasonError {
	return base.NewBaseOperationProcessReasonError(utils.ErrStringWrap(fmt.Sprintf("%s already exists, %s", name, k), e))
}

func ErrInvalid(t any, e error) base.BaseOperationProcessReasonError {
	return base.NewBaseOperationProcessReasonError(utils.ErrStringWrap(utils.ErrStringInvalid(t), e))
}
