package model

import (
	"strings"
)

//NewValidationError creates new validation error
func NewValidationError(messages []string) *ValidationError {
	return &ValidationError{messages: messages}
}

//ValidationError represents model validation error
type ValidationError struct {
	messages []string
}

//Error implementation error interface
func (e *ValidationError) Error() string {
	return strings.Join(e.messages, "\n")
}
