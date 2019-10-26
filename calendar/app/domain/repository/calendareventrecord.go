package repository

import (
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

//CalendarEventRecord repository record for Calendar Event
type CalendarEventRecord struct {
	ID          uint32    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Time        time.Time `json:"time"`
	CalendarID  uint32    `json:"calendarId"`
	Created     time.Time `json:"created"`
	LastUpdated time.Time `json:"lastUpdated"`
}

//CopyFromModel copy fields from model to repository record
func (r *CalendarEventRecord) CopyFromModel(m *model.CalendarEvent) {
	r.Title = m.Title
	r.Description = m.Description
	r.Location = m.Location
	r.Time = m.Time
	r.CalendarID = m.CalendarID
}

//CopyToModel copy fields to model from this record
func (r *CalendarEventRecord) CopyToModel(m *model.CalendarEvent) {
	m.Title = r.Title
	m.Description = r.Description
	m.Location = r.Location
	m.Time = r.Time
	m.CalendarID = r.CalendarID
}

//NewCalendarEventRecord create new record from calendar event model
func NewCalendarEventRecord(m *model.CalendarEvent) *CalendarEventRecord {
	r := &CalendarEventRecord{ID: m.ID}
	r.CopyFromModel(m)
	return r
}

//NewCalendarEventModel creates new model from calendar event repository record
func NewCalendarEventModel(r *CalendarEventRecord) *model.CalendarEvent {
	m := &model.CalendarEvent{ID: r.ID}
	r.CopyToModel(m)
	return m
}
