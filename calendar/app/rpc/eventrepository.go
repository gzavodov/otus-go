package rpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gzavodov/otus-go/calendar/app/domain/model"
	"github.com/gzavodov/otus-go/calendar/app/domain/repository"

	"google.golang.org/grpc"
)

//EventRepository RPC implementation of EventRepository interface
type EventRepository struct {
	ctx           context.Context
	serverAddress string
	conn          *grpc.ClientConn
	client        EventServiceClient
}

//NewEventRepository creates new RPC EventRepository
func NewEventRepository(ctx context.Context, serverAddress string) *EventRepository {
	if ctx == nil {
		ctx = context.Background()
	}
	return &EventRepository{ctx: ctx, serverAddress: serverAddress}
}

//Connect tries to connect to RPC server
func (r *EventRepository) Connect() (EventServiceClient, error) {
	if r.client != nil && r.conn != nil {
		return r.client, nil
	}

	var err error
	r.conn, err = grpc.Dial(r.serverAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	r.client = NewEventServiceClient(r.conn)
	return r.client, nil
}

//Disconnect closes connection with RPC server
func (r *EventRepository) Disconnect() {
	if r.client != nil {
		r.client = nil
	}

	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}
}

//Create add Calendar Event in repository
//If succseed ID field updated
func (r *EventRepository) Create(m *model.Event) error {
	if m == nil {
		return errors.New("first parameter must be not null pointer to event")
	}

	client, err := r.Connect()
	if err != nil {
		return err
	}
	defer r.Disconnect()

	event, err := CreateEventProxy(m)
	if err != nil {
		return err
	}

	event, err = client.Create(r.ctx, event)
	if err != nil {
		return err
	}

	m.ID = event.ID
	return nil
}

//Read get Calendar Event from repository by ID
func (r *EventRepository) Read(ID int64) (*model.Event, error) {
	if ID <= 0 {
		return nil, repository.NewError(repository.ErrorInvalidArgument, fmt.Sprintf("parameter 'ID' is invalid: %d", ID))
	}

	client, err := r.Connect()
	if err != nil {
		return nil, err
	}
	defer r.Disconnect()

	event, err := client.Read(r.ctx, &EventIdentifier{Value: ID})
	if err != nil {
		return nil, err
	}
	return CreateEventModel(event)
}

//ReadAll get all Calendar Events from repository
func (r *EventRepository) ReadAll() ([]*model.Event, error) {
	return nil, errors.New("method ReadAll is not supported in RPC context")
}

//ReadList get Calendar Events by interval specified by from and to params
func (r *EventRepository) ReadList(userID int64, from time.Time, to time.Time) ([]*model.Event, error) {
	client, err := r.Connect()
	if err != nil {
		return nil, err
	}
	defer r.Disconnect()

	query := EventListQuery{UserID: userID}

	fromTimestamp, err := ptypes.TimestampProto(from)
	if err != nil {
		return nil, err
	}
	query.From = fromTimestamp

	toTimestamp, err := ptypes.TimestampProto(to)
	if err != nil {
		return nil, err
	}
	query.To = toTimestamp

	reply, err := client.ReadList(r.ctx, &query)
	if err != nil {
		return nil, err
	}

	list := make([]*model.Event, 0, len(reply.Items))
	for _, event := range reply.Items {
		m, err := CreateEventModel(event)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

//ReadNotificationList get calendar events for notification
func (r *EventRepository) ReadNotificationList(userID int64, from time.Time) ([]*model.Event, error) {
	client, err := r.Connect()
	if err != nil {
		return nil, err
	}
	defer r.Disconnect()

	query := EventListQuery{UserID: userID}

	fromTimestamp, err := ptypes.TimestampProto(from)
	if err != nil {
		return nil, err
	}
	query.From = fromTimestamp

	reply, err := client.ReadNotificationList(r.ctx, &query)
	if err != nil {
		return nil, err
	}

	list := make([]*model.Event, 0, len(reply.Items))
	for _, event := range reply.Items {
		m, err := CreateEventModel(event)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

//IsExists check if repository contains Calendar event with specified ID
func (r *EventRepository) IsExists(ID int64) (bool, error) {
	return false, errors.New("method IsExists is not supported in RPC context")
}

//Update modifies Calendar Event in repository
func (r *EventRepository) Update(m *model.Event) error {
	if m == nil {
		return repository.NewError(repository.ErrorInvalidArgument, "first parameter must be not null pointer to event")
	}

	ID := m.ID
	if ID <= 0 {
		return repository.NewError(repository.ErrorInvalidArgument, fmt.Sprintf("model ID is invalid: %d", ID))
	}

	client, err := r.Connect()
	if err != nil {
		return err
	}
	defer r.Disconnect()

	event, err := CreateEventProxy(m)
	if err != nil {
		return err
	}

	_, err = client.Update(r.ctx, event)
	return err
}

//Delete removes Calendar Event from repository by ID
func (r *EventRepository) Delete(ID int64) error {
	if ID <= 0 {
		return repository.NewError(repository.ErrorInvalidArgument, fmt.Sprintf("parameter 'ID' is invalid: %d", ID))
	}

	client, err := r.Connect()
	if err != nil {
		return err
	}
	defer r.Disconnect()

	_, err = client.Delete(r.ctx, &EventIdentifier{Value: ID})
	return err
}

//GetTotalCount returns overall amouunt of calendar events in repository
func (r *EventRepository) GetTotalCount() (int64, error) {
	return 0, errors.New("method GetTotalCount is not supported in RPC context")
}
