package test

import (
	"errors"
	"fmt"
)

func NewObjectNotFoundError() error {
	return errors.New("object can not be found")
}

func NewObjectNotDeletedError() error {
	return errors.New("object can not be deleted")
}

func NewObjectNotMatchedError(expected, received interface{}) error {
	return fmt.Errorf("object before saving in repository is not equal to object after reading from repository; expected: %v, received: %v", expected, received)
}

func NewUnexpectedResponseCodeError(expected, received int) error {
	return fmt.Errorf("unexpected response code: got %v want %v", expected, received)
}
