package repository

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

//InMemoryCalendarEventRepository thread safe in-memory implementation of CalendarEventRepository interface
type InMemoryCalendarEventRepository struct {
	mu           *sync.RWMutex
	records      map[uint32]*CalendarEventRecord
	lastRecordID uint32
}

//NewInMemoryCalendarEventRepository creates new in-memory CalendarEventRepository
func NewInMemoryCalendarEventRepository() *InMemoryCalendarEventRepository {
	return &InMemoryCalendarEventRepository{
		mu:           &sync.RWMutex{},
		records:      make(map[uint32]*CalendarEventRecord),
		lastRecordID: 0,
	}
}

//Create add Calendar Event in repository
//If succseed ID field updated
func (r *InMemoryCalendarEventRepository) Create(m *model.CalendarEvent) error {
	if m == nil {
		return errors.New("parameter 'm' must be not null pointer")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	record := NewCalendarEventRecord(m)
	record.Created = time.Now()
	record.LastUpdated = record.Created

	r.lastRecordID++
	record.ID = r.lastRecordID
	m.ID = record.ID

	r.records[record.ID] = record

	return nil
}

//Read get Calendar Event from repository by ID
func (r *InMemoryCalendarEventRepository) Read(ID uint32) (*model.CalendarEvent, error) {
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
func (r *InMemoryCalendarEventRepository) ReadAll() []*model.CalendarEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*model.CalendarEvent, 0, len(r.records))
	for _, record := range r.records {
		list = append(list, NewCalendarEventModel(record))
	}

	sort.SliceStable(
		list,
		func(i, j int) bool { return list[i].ID < list[j].ID },
	)

	return list
}

//IsExists check if repository contains Calendar event with specified ID
func (r *InMemoryCalendarEventRepository) IsExists(ID uint32) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, isFound := r.records[ID]
	return isFound
}

//Update modifies Calendar Event in repository
func (r *InMemoryCalendarEventRepository) Update(m *model.CalendarEvent) error {
	if m == nil {
		return errors.New("parameter 'm' must be not null pointer")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	ID := m.ID
	if ID <= 0 {
		return fmt.Errorf("model ID is invalid: %d", ID)
	}

	record, isFound := r.records[ID]
	if !isFound {
		return fmt.Errorf("could not find record with ID: %d", ID)
	}

	record.CopyFromModel(m)
	record.LastUpdated = time.Now()

	return nil
}

//Delete removes Calendar Event from repository by ID
func (r *InMemoryCalendarEventRepository) Delete(ID uint32) error {
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
func (r *InMemoryCalendarEventRepository) GetTotalCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.records)
}

//Purge removes all Calendar records from repository
func (r *InMemoryCalendarEventRepository) Purge() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = make(map[uint32]*CalendarEventRecord)
	r.lastRecordID = 0

	return nil
}
