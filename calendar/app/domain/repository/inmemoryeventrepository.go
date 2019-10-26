package repository

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

type InMemoryEventRepository struct {
	mu           *sync.RWMutex
	records      map[uint32]*CalendarRecord
	lastRecordID uint32
}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		mu:           &sync.RWMutex{},
		records:      make(map[uint32]*CalendarRecord),
		lastRecordID: 0,
	}
}

func (r *InMemoryEventRepository) Create(m *model.CalendarEvent) error {
	if m == nil {
		return errors.New("parameter 'm' must be not null pointer")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	record := NewRecord(m)

	record.Created = time.Now()
	record.LastUpdated = record.Created

	r.lastRecordID++
	record.ID = r.lastRecordID
	m.ID = record.ID
	r.records[record.ID] = record
	return nil
}

func (r *InMemoryEventRepository) Read(ID uint32) (model.CalendarEvent, error) {
	if ID <= 0 {
		return nil, fmt.Errorf("parameter 'ID' is invalid: %d", ID)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	record, isFound := r.records[ID]
	if !isFound {
		return nil, fmt.Errorf("could not find record with ID: %d", ID)
	}

	return NewModel(record), nil
}

func (r *InMemoryEventRepository) ReadAll() []*model.CalendarEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*model.CalendarEvent, 0, len(r.records))
	for _, record := range r.records {
		list = append(list, NewModel(record))
	}
	return list
}

func (r *InMemoryEventRepository) Update(m model.CalendarEvent) error {
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
