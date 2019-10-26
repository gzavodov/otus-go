package repository

import (
	"testing"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

func TestInMemoryRepository(t *testing.T) {
	repo := &InMemoryEventRepository{}

	model := model.CalendarEvent{Title: "Test Event #1", Time: time.Parse(time.RFC3339, "2019-10-01T12:00:00Z02:00")}
	repo.Create(model)
}
