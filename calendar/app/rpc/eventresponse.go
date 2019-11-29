package rpc

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/gzavodov/otus-go/calendar/app/domain/model"
	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//EventResponseBase base structure for response
type EventResponseBase struct {
	Handler *EventHandler
}

//LogError writes error in log
func (r *EventResponseBase) LogError(errorName string, err error) {
	r.Handler.LogError(errorName, err)
}

//CreateEventProxy creates new GRPC event proxy from model
func (r *EventResponseBase) CreateEventProxy(eventModel *model.Event) (*Event, error) {
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

	eventProxy.UserID = eventModel.UserID
	eventProxy.CalendarID = eventModel.CalendarID

	return eventProxy, nil
}

//PrepareReplyError returns GRPC error
func (r *EventResponseBase) PrepareReplyError(err error) error {
	if err == nil {
		return nil
	}

	code := codes.Internal
	if _, ok := err.(*model.ValidationError); ok {
		code = codes.FailedPrecondition
	} else if repositoryErr, ok := err.(*repository.Error); ok {
		switch repositoryErr.GetCode() {
		case repository.ErrorNotFound:
			code = codes.NotFound
		case repository.ErrorInvalidArgument:
			code = codes.InvalidArgument
		case repository.ErrorInvalidObject:
			code = codes.FailedPrecondition
		}
	}

	return status.Error(code, err.Error())
}

//EventResponse structure for response with event
type EventResponse struct {
	IncomingProxy *Event

	EventResponseBase
}

//NewEventResponse constructs new EventResponse
func NewEventResponse(handler *EventHandler, incomingProxy *Event) *EventResponse {
	return &EventResponse{EventResponseBase: EventResponseBase{Handler: handler}, IncomingProxy: incomingProxy}
}

//LogAndReply writes error in log and returns GPRC response data
func (r *EventResponse) LogAndReply(model *model.Event, errorName string, err error) (*Event, error) {
	if err != nil {
		r.LogError(errorName, err)
		return r.IncomingProxy, r.PrepareReplyError(err)
	}

	var outgoingProxy *Event
	if model != nil {
		outgoingProxy, err = r.CreateEventProxy(model)
		if err != nil {
			r.LogError(ErrorCategoryExternalization, err)
		}
	}

	if outgoingProxy == nil {
		outgoingProxy = r.IncomingProxy
	}

	return outgoingProxy, r.PrepareReplyError(err)
}

//EventIdentifierResponse structure for response with event ID
type EventIdentifierResponse struct {
	IncomingProxy *EventIdentifier

	EventResponseBase
}

//NewEventIdentifierResponse constructs new EventIdentifierResponse
func NewEventIdentifierResponse(handler *EventHandler, incomingProxy *EventIdentifier) *EventIdentifierResponse {
	return &EventIdentifierResponse{EventResponseBase: EventResponseBase{Handler: handler}, IncomingProxy: incomingProxy}
}

//LogAndReply writes error in log and returns GPRC response data
func (r *EventIdentifierResponse) LogAndReply(errorName string, err error) (*EventIdentifier, error) {
	if err != nil {
		r.LogError(errorName, err)
	}
	return r.IncomingProxy, r.PrepareReplyError(err)
}

//EventListResponse structure for response with list of events
type EventListResponse struct {
	EventResponseBase
}

//NewEventListResponse constructs new EventListResponse
func NewEventListResponse(handler *EventHandler) *EventListResponse {
	return &EventListResponse{EventResponseBase: EventResponseBase{Handler: handler}}
}

//LogAndReply writes error in log and returns GPRC response data
func (r *EventListResponse) LogAndReply(models []*model.Event, errorName string, err error) (*EventListReply, error) {
	reply := &EventListReply{}

	if err != nil {
		r.LogError(errorName, err)
		return reply, err
	}

	if len(models) > 0 {
		proxes := make([]*Event, 0, len(models))
		for _, model := range models {
			proxy, err := r.CreateEventProxy(model)
			if err != nil {
				r.LogError(ErrorCategoryExternalization, err)
				return reply, err
			}
			proxes = append(proxes, proxy)
		}
		reply.Items = proxes
	}

	return reply, nil
}
