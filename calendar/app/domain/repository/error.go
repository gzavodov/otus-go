package repository

import (
	"errors"
	"fmt"
)

//Repository Error Type
const (
	ErrorNone            = 0
	ErrorGeneral         = 1
	ErrorNotFound        = 2
	ErrorInvalidArgument = 3
	ErrorInvalidObject   = 4

	ErrorInvalidConfiguration = 5
	ErrorFailedToConnect      = 6

	ErrorDataCreationFailure     = 7
	ErrorDataRetrievalFailure    = 8
	ErrorDataModificationFailure = 9
	ErrorDataDeletionFailure     = 10
)

//NewError creates new repository error
func NewError(code int, msg string) *Error {
	return &Error{code: code, innerErr: errors.New(msg)}
}

//NewErrorf creates new repository error with formatted message
func NewErrorf(code int, msg string, args ...interface{}) *Error {
	return &Error{code: code, innerErr: fmt.Errorf(msg, args...)}
}

//WrapError creates a new repository with wrapped error
func WrapError(code int, err error) error {
	return &Error{code: code, innerErr: err}
}

//WrapErrorf creates a new repository with wrapped error and formatted message
func WrapErrorf(code int, err error, msg string, args ...interface{}) error {
	return &Error{code: code, innerErr: fmt.Errorf("%s: underlying error (%w)", fmt.Sprintf(msg, args...), err)}
}

//Error represents custom repository error
type Error struct {
	code     int
	innerErr error
}

//GetCode retuns repository error code
func (e *Error) GetCode() int { return e.code }

//Unwrap implementation of Unwrap interface
func (e *Error) Unwrap() error { return e.innerErr }

//Error implementation error interface
func (e *Error) Error() string {
	if e.innerErr == nil {
		return ""
	}

	return e.innerErr.Error()
}
