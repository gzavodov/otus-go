package inmemory

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
)

//EventRepository thread safe in-memory implementation of EventRepository interface
type EventRepository struct {
	mu           sync.RWMutex
	records      map[int64]*repository.EventRecord
	lastRecordID int64
}

//NewEventRepository creates new in-memory EventRepository
func NewEventRepository() *EventRepository {
	return &EventRepository{
		mu:      sync.RWMutex{},
		records: make(map[int64]*repository.EventRecord),
	}
}

//Create add Calendar Event in repository
//If succseed ID field updated
func (r *EventRepository) Create(m *model.Event) error {
	if m == nil {
		return errors.New("first parameter must be not null pointer to event")
	}

	record := repository.NewCalendarEventRecord(m)
	record.Created = time.Now()
	record.LastUpdated = record.Created

	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastRecordID++
	record.ID = r.lastRecordID
	m.ID = record.ID

	r.records[record.ID] = record

	return nil
}

//Read get Calendar Event from repository by ID
func (r *EventRepository) Read(ID int64) (*model.Event, error) {
	if ID <= 0 {
		return nil, repository.NewError(repository.ErrorInvalidArgument, fmt.Sprintf("parameter 'ID' is invalid: %d", ID))
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	record, isFound := r.records[ID]
	if !isFound {
		return nil, repository.NewError(repository.ErrorNotFound, fmt.Sprintf("could not find record with ID: %d", ID))
	}

	return repository.NewCalendarEventModel(record), nil
}

//ReadAll get all Calendar Events from repository
func (r *EventRepository) ReadAll() ([]*model.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*model.Event, 0, len(r.records))
	for _, record := range r.records {
		list = append(list, repository.NewCalendarEventModel(record))
	}

	sort.SliceStable(
		list,
		func(i, j int) bool { return list[i].ID < list[j].ID },
	)

	return list, nil
}

//ReadList get Calendar Events by interval specified by from and to params
func (r *EventRepository) ReadList(userID int64, from time.Time, to time.Time) ([]*model.Event, error) {
	list := make([]*model.Event, 0, len(r.records))
	for _, record := range r.records {
		if userID > 0 && record.UserID != userID {
			continue
		}
		if (record.StartTime.Equal(from) || record.StartTime.After(from)) && (record.EndTime.Equal(to) || record.EndTime.Before(to)) {
			list = append(list, repository.NewCalendarEventModel(record))
		}
	}

	sort.SliceStable(
		list,
		func(i, j int) bool { return list[i].ID < list[j].ID },
	)

	return list, nil
}

//ReadNotificationList get calendar events for notification
func (r *EventRepository) ReadNotificationList(userID int64, from time.Time) ([]*model.Event, error) {
	from = from.UTC()

	list := make([]*model.Event, 0, len(r.records))
	for _, record := range r.records {
		if userID > 0 && record.UserID != userID {
			continue
		}

		currentNotifyTime := record.StartTime.Add(time.Duration(-1 * int64(record.NotifyBefore)))
		isMatched := (from.Before(record.StartTime) || from.Equal(record.StartTime)) &&
			(from.After(currentNotifyTime) || from.Equal(currentNotifyTime))
		if isMatched {
			list = append(list, repository.NewCalendarEventModel(record))
		}
	}

	sort.SliceStable(
		list,
		func(i, j int) bool { return list[i].ID < list[j].ID },
	)

	return list, nil
}

//IsExists check if repository contains Calendar event with specified ID
func (r *EventRepository) IsExists(ID int64) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, isFound := r.records[ID]
	return isFound, nil
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

	r.mu.Lock()
	defer r.mu.Unlock()

	record, isFound := r.records[ID]
	if !isFound {
		return repository.NewError(repository.ErrorNotFound, fmt.Sprintf("could not find record with ID: %d", ID))
	}

	record.CopyFromModel(m)
	record.LastUpdated = time.Now()

	return nil
}

//Delete removes Calendar Event from repository by ID
func (r *EventRepository) Delete(ID int64) error {
	if ID <= 0 {
		return repository.NewError(repository.ErrorInvalidArgument, fmt.Sprintf("parameter 'ID' is invalid: %d", ID))
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, isFound := r.records[ID]; !isFound {
		return repository.NewError(repository.ErrorNotFound, fmt.Sprintf("could not find record with ID: %d", ID))
	}

	delete(r.records, ID)

	return nil
}

//GetTotalCount returns overall amouunt of calendar events in repository
func (r *EventRepository) GetTotalCount() (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return int64(len(r.records)), nil
}

//Purge removes all Calendar records from repository
func (r *EventRepository) purge() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = make(map[int64]*repository.EventRecord)
	r.lastRecordID = 0

	return nil
}
