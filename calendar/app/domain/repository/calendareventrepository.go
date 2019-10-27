package repository

import (
	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

//CalendarEventRepository Storage interface for Calendar Event
type CalendarEventRepository interface {
	Create(*model.CalendarEvent) error
	Read(uint32) (*model.CalendarEvent, error)
	ReadAll() []*model.CalendarEvent
	IsExists(uint32) bool
	Update(*model.CalendarEvent) error
	Delete(uint32) error
	GetTotalCount() int
	Purge() error
}
