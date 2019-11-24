package model

import (
	"fmt"
)

//EventValidationResult calendar event validation result
type EventValidationResult struct {
	messages []string
}

func (v *EventValidationResult) addMessage(msg string) {
	v.messages = append(v.messages, msg)
}

func (v *EventValidationResult) addRequiredFieldMessage(name string) {
	v.addMessage(fmt.Sprintf("field \"%s\" is required", name))
}

//GetMessages return validation error messages
func (v *EventValidationResult) GetMessages() []string {
	return v.messages
}

//EventValidator calendar event standard validator
type EventValidator struct {
	Event *Event
}

//Validate performs calendar event validation
func (v *EventValidator) Validate() *EventValidationResult {
	result := &EventValidationResult{}
	if len(v.Event.Title) == 0 {
		result.addRequiredFieldMessage("Title")
	}

	if v.Event.UserID <= 0 {
		result.addRequiredFieldMessage("UserID")
	}

	if v.Event.StartTime.IsZero() {
		result.addRequiredFieldMessage("StartTime")
	}

	if v.Event.EndTime.IsZero() {
		result.addRequiredFieldMessage("EndTime")
	}

	if !v.Event.StartTime.Before(v.Event.EndTime) {
		result.addMessage("StartTime must be before EndTime")
	}

	return result
}
