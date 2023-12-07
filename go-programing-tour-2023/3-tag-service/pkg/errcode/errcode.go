package errcode

import (
	"fmt"
)

type Error struct {
	code int
	msg  string
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Error() string {
	return fmt.Sprintf("error code: %d, error message: %s", e.Code(), e.Msg())
}

func NewError(code int, msg string) *Error {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("error code %d already exist", code))
	}
	_codes[code] = msg
	return &Error{code: code, msg: msg}
}

var _codes = map[int]string{}
