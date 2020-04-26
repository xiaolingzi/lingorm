package common

import (
	"errors"
	"fmt"
	"runtime/debug"
)

type Error struct {
	Message string
	Stack   string
}

func NewError() *Error {
	return &Error{}
}

func (h *Error) Throw(err interface{}) {
	panic(err)
}

func (h *Error) Defer(err *error) {
	if r := recover(); r != nil {
		var tempErr error
		switch x := r.(type) {
		case string:
			tempErr = errors.New(x)
		case error:
			tempErr = x
		default:
			tempErr = errors.New(fmt.Sprintf("%s", r))
		}
		h.Message = tempErr.Error()
		h.Stack = string(debug.Stack())
		*err = h
	}
}

func (h *Error) Error() string {
	return h.Message + "\n" + h.Stack
}
