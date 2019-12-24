package repository

import (
	"time"

	"github.com/gzavodov/otus-go/calendar/model"
)

//EventRepository Storage interface for Calendar Event
type EventRepository interface {
	Create(*model.Event) error
	Read(int64) (*model.Event, error)
	ReadAll() ([]*model.Event, error)
	ReadList(userID int64, from time.Time, to time.Time) ([]*model.Event, error)
	ReadNotificationList(userID int64, from time.Time) ([]*model.Event, error)
	IsExists(int64) (bool, error)
	Update(*model.Event) error
	Delete(int64) error
	GetTotalCount() (int64, error)
}
