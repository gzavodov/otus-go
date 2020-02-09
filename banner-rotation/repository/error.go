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

	ErrorRecordCreationFailure = 7
	ErrorRecordReadingFailure  = 8
	ErrorRecordUpdatingFailure = 9
	ErrorRecordDeletionFailure = 10
)

//NewError creates new repository error
func NewError(code int, msg string) *Error {
	return &Error{code: code, innerErr: errors.New(msg)}
}

//NewErrorf creates new repository error with formatted message
func NewErrorf(code int, msg string, args ...interface{}) *Error {
	return &Error{code: code, innerErr: fmt.Errorf(msg, args...)}
}

//NewInvalidArgumentError creates new Invalid Argument Error with specified formatted message
func NewInvalidArgumentError(msg string, args ...interface{}) *Error {
	return NewErrorf(ErrorInvalidArgument, msg, args...)
}

//NewNotFoundError creates new Not Found Error with specified formatted message
func NewNotFoundError(msg string, args ...interface{}) *Error {
	return NewErrorf(ErrorNotFound, msg, args...)
}

//NewCreationError creates new Record Creation Error with specified inner error and message
func NewCreationError(err error, msg string) *Error {
	return WrapErrorf(ErrorRecordCreationFailure, err, msg)
}

//NewReadingError creates new Record Reading Error with specified inner error and formatted message
func NewReadingError(err error, msg string, args ...interface{}) *Error {
	return WrapErrorf(ErrorRecordReadingFailure, err, msg, args...)
}

//NewUpdatingError creates new Record Updating Error with specified inner error and message
func NewUpdatingError(err error, msg string, args ...interface{}) *Error {
	return WrapErrorf(ErrorRecordUpdatingFailure, err, msg, args...)
}

//NewDeletionError creates new Record Deletion Error with specified inner error and message
func NewDeletionError(err error, msg string, args ...interface{}) *Error {
	return WrapErrorf(ErrorRecordDeletionFailure, err, msg, args...)
}

//WrapError creates a new repository with wrapped error
func WrapError(code int, err error) error {
	return &Error{code: code, innerErr: err}
}

//WrapErrorf creates a new repository with wrapped error and formatted message
func WrapErrorf(code int, err error, msg string, args ...interface{}) *Error {
	if err != nil {
		return &Error{code: code, innerErr: fmt.Errorf("%s: underlying error (%w)", fmt.Sprintf(msg, args...), err)}
	}
	return &Error{code: code, innerErr: fmt.Errorf(msg, args...)}
}

func IsNotFoundError(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}

	return e.GetCode() == ErrorNotFound
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
