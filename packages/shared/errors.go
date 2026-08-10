package shared

import (
	"errors"
	"fmt"
)

// Code identifies a stable error class for APIs and logs.
type Code string

const (
	CodeInvalid   Code = "invalid"
	CodeNotFound  Code = "not_found"
	CodeConflict  Code = "conflict"
	CodeForbidden Code = "forbidden"
	CodeInternal  Code = "internal"
)

// Error is a domain/application error with a stable code.
type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func NewError(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(code Code, message string, err error) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

func AsError(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}
