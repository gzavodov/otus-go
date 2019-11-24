package repository

import (
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

//EventRecord repository record for Calendar Event
type EventRecord struct {
	ID          uint32    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	UserID      uint32    `json:"userId"`
	CalendarID  uint32    `json:"calendarId"`
	Created     time.Time `json:"created"`
	LastUpdated time.Time `json:"lastUpdated"`
}

//CopyFromModel copy fields from model to repository record
func (r *EventRecord) CopyFromModel(m *model.Event) {
	r.Title = m.Title
	r.Description = m.Description
	r.Location = m.Location
	r.StartTime = m.StartTime
	r.EndTime = m.EndTime
	r.UserID = m.UserID
	r.CalendarID = m.CalendarID
}

//CopyToModel copy fields to model from this record
func (r *EventRecord) CopyToModel(m *model.Event) {
	m.Title = r.Title
	m.Description = r.Description
	m.Location = r.Location
	m.StartTime = r.StartTime
	m.EndTime = r.EndTime
	m.UserID = r.UserID
	m.CalendarID = r.CalendarID
}

//NewCalendarEventRecord create new record from calendar event model
func NewCalendarEventRecord(m *model.Event) *EventRecord {
	r := &EventRecord{ID: m.ID}
	r.CopyFromModel(m)
	return r
}

//NewCalendarEventModel creates new model from calendar event repository record
func NewCalendarEventModel(r *EventRecord) *model.Event {
	m := &model.Event{ID: r.ID}
	r.CopyToModel(m)
	return m
}
