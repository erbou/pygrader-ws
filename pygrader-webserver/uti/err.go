package uti

import (
	"fmt"
)

type ErrorCode int

type Error struct {
	error
	Code ErrorCode
}

const (
	ERR_SYSTEM_ERROR ErrorCode = 1 + iota
)

func (e Error) Error() string {
	return fmt.Sprintf("error[%v]: %v", e.Code, e.error.Error())
}

func (e Error) Unwrap() error {
	return e.error
}

func Errorf(code ErrorCode, format string, args ...interface{}) Error {
	return Error{Code: code, error: fmt.Errorf(format, args...)}
}
