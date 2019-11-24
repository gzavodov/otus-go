package repository

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

//InMemoryEventRepository thread safe in-memory implementation of EventRepository interface
type InMemoryEventRepository struct {
	mu           sync.RWMutex
	records      map[uint32]*EventRecord
	lastRecordID uint32
}

//NewInMemoryEventRepository creates new in-memory EventRepository
func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		mu:      sync.RWMutex{},
		records: make(map[uint32]*EventRecord),
	}
}

//Create add Calendar Event in repository
//If succseed ID field updated
func (r *InMemoryEventRepository) Create(m *model.Event) error {
	if m == nil {
		return errors.New("parameter 'm' must be not null pointer")
	}

	validator := model.EventValidator{Event: m}
	errorMessages := validator.Validate().GetMessages()
	if len(errorMessages) > 0 {
		return errors.New(strings.Join(errorMessages, "\n"))
	}

	record := NewCalendarEventRecord(m)
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
func (r *InMemoryEventRepository) Read(ID uint32) (*model.Event, error) {
	if ID <= 0 {
		return nil, fmt.Errorf("parameter 'ID' is invalid: %d", ID)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	record, isFound := r.records[ID]
	if !isFound {
		return nil, fmt.Errorf("could not find record with ID: %d", ID)
	}

	return NewCalendarEventModel(record), nil
}

//ReadAll get all Calendar Events from repository
func (r *InMemoryEventRepository) ReadAll() []*model.Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*model.Event, 0, len(r.records))
	for _, record := range r.records {
		list = append(list, NewCalendarEventModel(record))
	}

	sort.SliceStable(
		list,
		func(i, j int) bool { return list[i].ID < list[j].ID },
	)

	return list
}

//ReadList get Calendar Events by interval specified by from and to params
func (r *InMemoryEventRepository) ReadList(userID uint32, from time.Time, to time.Time) ([]*model.Event, error) {
	list := make([]*model.Event, 0, len(r.records))
	for _, record := range r.records {
		if userID > 0 && record.UserID != userID {
			continue
		}
		if (record.StartTime.Equal(from) || record.StartTime.After(from)) && (record.EndTime.Equal(to) || record.EndTime.Before(to)) {
			list = append(list, NewCalendarEventModel(record))
		}
	}

	sort.SliceStable(
		list,
		func(i, j int) bool { return list[i].ID < list[j].ID },
	)

	return list, nil
}

//IsExists check if repository contains Calendar event with specified ID
func (r *InMemoryEventRepository) IsExists(ID uint32) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, isFound := r.records[ID]
	return isFound
}

//Update modifies Calendar Event in repository
func (r *InMemoryEventRepository) Update(m *model.Event) error {
	if m == nil {
		return errors.New("parameter 'm' must be not null pointer")
	}

	ID := m.ID
	if ID <= 0 {
		return fmt.Errorf("model ID is invalid: %d", ID)
	}

	validator := model.EventValidator{Event: m}
	validationMessages := validator.Validate().GetMessages()
	if len(validationMessages) > 0 {
		return errors.New(strings.Join(validationMessages, "\n"))
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	record, isFound := r.records[ID]
	if !isFound {
		return fmt.Errorf("could not find record with ID: %d", ID)
	}

	record.CopyFromModel(m)
	record.LastUpdated = time.Now()

	return nil
}

//Delete removes Calendar Event from repository by ID
func (r *InMemoryEventRepository) Delete(ID uint32) error {
	if ID <= 0 {
		return fmt.Errorf("parameter 'ID' is invalid: %d", ID)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, isFound := r.records[ID]; !isFound {
		return fmt.Errorf("could not find record with ID: %d", ID)
	}

	delete(r.records, ID)

	return nil
}

//GetTotalCount returns overall amouunt of calendar events in repository
func (r *InMemoryEventRepository) GetTotalCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.records)
}

//Purge removes all Calendar records from repository
func (r *InMemoryEventRepository) Purge() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = make(map[uint32]*EventRecord)
	r.lastRecordID = 0

	return nil
}
