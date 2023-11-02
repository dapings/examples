package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	// 错误码
	code int `json:"code"`
	// 错误消息
	msg string `json:"msg"`
	// 详细信息
	details []string `json:"details"`
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Details() []string {
	return e.details
}

func (e *Error) Error() string {
	return fmt.Sprintf("errCode: %d, errMsg: %s", e.Code(), e.Msg())
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("err code %d exist", code))
	}
	codes[code] = msg
	return &Error{
		code: code,
		msg:  msg,
	}
}

func (e *Error) WithDetails(details ...string) *Error {
	newErr := *e
	newErr.details = []string{}
	for _, d := range details {
		newErr.details = append(newErr.details, d)
	}
	return &newErr
}

func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ServerError.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedTokenError.Code():
		fallthrough
	case UnauthorizedTokenGenerate.Code():
		fallthrough
	case UnauthorizedTokenTimeout.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code():
		return http.StatusTooManyRequests
	}

	return http.StatusInternalServerError
}
