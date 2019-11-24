package repository

import (
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

//EventRepository Storage interface for Calendar Event
type EventRepository interface {
	Create(*model.Event) error
	Read(uint32) (*model.Event, error)
	ReadAll() []*model.Event
	ReadList(userID uint32, from time.Time, to time.Time) ([]*model.Event, error)
	IsExists(uint32) bool
	Update(*model.Event) error
	Delete(uint32) error
	GetTotalCount() int
	Purge() error
}
