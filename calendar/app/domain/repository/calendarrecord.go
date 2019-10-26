package repository

import (
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

type CalendarRecord struct {
	ID          uint32
	Title       string
	Description string
	Location    string
	Time        time.Time
	CalendarID  uint32
	Created     time.Time
	LastUpdated time.Time
}

func (r *CalendarRecord) CopyFromModel(m *model.CalendarEvent) {
	r.Title = m.Title
	r.Description = m.Description
	r.Location = m.Location
	r.Time = m.Time
	r.CalendarID = m.CalendarID
}

func (r *CalendarRecord) CopyToModel(m *model.CalendarEvent) {
	m.Title = r.Title
	m.Description = r.Description
	m.Location = r.Location
	m.Time = r.Time
	m.CalendarID = r.CalendarID
}

func NewRecord(m *model.CalendarEvent) *CalendarRecord {
	r := &CalendarRecord{ID: m.ID}
	r.CopyFromModel(m)
	return r
}

func NewModel(r *CalendarRecord) *model.CalendarEvent {
	m := &model.CalendarEvent{ID: r.ID}
	r.CopyToModel(m)
	return m
}
