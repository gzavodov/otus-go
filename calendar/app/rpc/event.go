package rpc

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

//CreateEventProxy creates new GRPC event proxy from model
func CreateEventProxy(eventModel *model.Event) (*Event, error) {
	if eventModel == nil {
		return &Event{}, nil
	}

	eventProxy := &Event{
		ID:          eventModel.ID,
		Title:       eventModel.Title,
		Description: eventModel.Description,
		Location:    eventModel.Location,
	}

	startTime, err := ptypes.TimestampProto(eventModel.StartTime)
	if err != nil {
		return nil, err
	}
	eventProxy.StartTime = startTime

	endTime, err := ptypes.TimestampProto(eventModel.EndTime)
	if err != nil {
		return nil, err
	}
	eventProxy.EndTime = endTime

	eventProxy.NotifyBefore = eventModel.NotifyBefore.Nanoseconds()

	eventProxy.UserID = eventModel.UserID
	eventProxy.CalendarID = eventModel.CalendarID

	return eventProxy, nil
}

//CreateEventModel creates new model from GRPC event proxy
func CreateEventModel(eventProxy *Event) (*model.Event, error) {
	if eventProxy == nil {
		return &model.Event{}, nil
	}

	eventModel := &model.Event{}
	eventModel.ID = eventProxy.ID
	eventModel.Title = eventProxy.Title
	eventModel.Description = eventProxy.Description
	eventModel.Location = eventProxy.Location

	startTime, err := ptypes.Timestamp(eventProxy.StartTime)
	if err != nil {
		return nil, err
	}
	eventModel.StartTime = startTime

	endTime, err := ptypes.Timestamp(eventProxy.EndTime)
	if err != nil {
		return nil, err
	}
	eventModel.EndTime = endTime

	eventModel.NotifyBefore = time.Duration(eventProxy.NotifyBefore)
	eventModel.UserID = eventProxy.UserID
	eventModel.CalendarID = eventProxy.CalendarID

	return eventModel, nil
}

//ValidateEventModel validates calendar event model
func ValidateEventModel(eventModel *model.Event) error {
	validator := model.EventValidator{Event: eventModel}
	return validator.Validate().Error()
}
