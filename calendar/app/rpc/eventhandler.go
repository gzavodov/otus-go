package rpc

import (
	context "context"

	"github.com/golang/protobuf/ptypes"
	"github.com/gzavodov/otus-go/calendar/app/domain/model"
	"github.com/gzavodov/otus-go/calendar/app/endpoint"
)

//EventHandler structure for event action handler
type EventHandler struct {
	endpoint.Handler
}

//Create creates new calendar event
func (h *EventHandler) Create(ctx context.Context, in *Event) (*Event, error) {
	response := NewEventResponse(h, in)

	eventModel, err := h.CreateEventModel(in)
	if err != nil {
		response.LogAndReply(nil, ErrorCategoryInternalization, err)
	}

	err = h.ValidateEventModel(eventModel)
	if err != nil {
		return response.LogAndReply(eventModel, ErrorCategoryValidation, err)
	}

	return response.LogAndReply(
		eventModel,
		ErrorCategoryRepository,
		h.Repo.Create(eventModel),
	)
}

//Read reads calendar event by ID
func (h *EventHandler) Read(ctx context.Context, in *EventIdentifier) (*Event, error) {
	response := NewEventResponse(h, nil)

	eventModel, err := h.Repo.Read(in.Value)
	return response.LogAndReply(
		eventModel,
		ErrorCategoryRepository,
		err,
	)
}

//Update updates calendar event
func (h *EventHandler) Update(ctx context.Context, in *Event) (*Event, error) {
	response := NewEventResponse(h, in)

	eventModel, err := h.CreateEventModel(in)
	if err != nil {
		response.LogAndReply(nil, ErrorCategoryInternalization, err)
	}

	err = h.ValidateEventModel(eventModel)
	if err != nil {
		return response.LogAndReply(eventModel, ErrorCategoryValidation, err)
	}

	return response.LogAndReply(
		eventModel,
		ErrorCategoryRepository,
		h.Repo.Update(eventModel),
	)
}

//Delete deletes calendar event by ID
func (h *EventHandler) Delete(ctx context.Context, in *EventIdentifier) (*EventIdentifier, error) {
	response := NewEventIdentifierResponse(h, in)
	return response.LogAndReply(
		ErrorCategoryRepository,
		h.Repo.Delete(in.Value),
	)
}

//ReadList returns list of calendar events selected by filter
func (h *EventHandler) ReadList(ctx context.Context, in *EventListQuery) (*EventListReply, error) {
	response := NewEventListResponse(h)

	from, err := ptypes.Timestamp(in.From)
	if err != nil {
		return response.LogAndReply(nil, ErrorCategoryInternalization, err)
	}

	to, err := ptypes.Timestamp(in.To)
	if err != nil {
		return response.LogAndReply(nil, ErrorCategoryInternalization, err)
	}

	eventModels, err := h.Repo.ReadList(in.UserID, from, to)
	return response.LogAndReply(eventModels, ErrorCategoryRepository, err)
}

//CreateEventModel creates new model from GRPC event proxy
func (h *EventHandler) CreateEventModel(eventProxy *Event) (*model.Event, error) {
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

	eventModel.UserID = eventProxy.UserID
	eventModel.CalendarID = eventProxy.CalendarID

	return eventModel, nil
}

//ValidateEventModel validates calendar event model
func (h *EventHandler) ValidateEventModel(eventModel *model.Event) error {
	validator := model.EventValidator{Event: eventModel}
	return validator.Validate().Error()
}
