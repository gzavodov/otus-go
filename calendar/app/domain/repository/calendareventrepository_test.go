package repository

import (
	"testing"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

func TestInMemoryRepository(t *testing.T) {
	repo := NewInMemoryCalendarEventRepository()

	t.Run("Checking of EventRepository::Create",
		func(t *testing.T) {
			initialEvent := &model.CalendarEvent{Title: "Test Event #1 (2019-10-01T12:00:00)", Time: time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC)}
			err := repo.Create(initialEvent)
			if err != nil {
				t.Fatal(err)
			}

			if initialEvent.ID <= 0 {
				t.Error("CREATION FAILED: Initail calendar event must have valid ID\n")
			}

			repo.Purge()
		})

	t.Run("Checking of EventRepository::Read",
		func(t *testing.T) {
			initialEvent := &model.CalendarEvent{Title: "Test Event #2 (2019-10-02T12:00:00)", Time: time.Date(2019, 10, 2, 12, 0, 0, 0, time.UTC)}
			err := repo.Create(initialEvent)
			if err != nil {
				t.Fatal(err)
			}

			resultEvent, err := repo.Read(initialEvent.ID)
			if err != nil {
				t.Fatal(err)
			}

			if *resultEvent != *initialEvent {
				t.Error("READING FAILED: Initail calendar event is not equals repository event\n")
			}

			repo.Purge()
		})

	t.Run("Checking of EventRepository::Update",
		func(t *testing.T) {
			initialEvent := &model.CalendarEvent{Title: "Test Event #3 (2019-10-03T12:00:00)", Time: time.Date(2019, 10, 3, 12, 0, 0, 0, time.UTC)}
			err := repo.Create(initialEvent)
			if err != nil {
				t.Fatal(err)
			}

			initialEvent.Description = "Test Description"
			err = repo.Update(initialEvent)
			if err != nil {
				t.Fatal(err)
			}

			resultEvent, err := repo.Read(initialEvent.ID)
			if err != nil {
				t.Fatal(err)
			}

			if *resultEvent != *initialEvent {
				t.Error("UPDATE FAILED: Initail calendar event is not equals repository event\n")
			}

			repo.Purge()
		})

	t.Run("Checking of EventRepository::Delete",
		func(t *testing.T) {
			initialEvent := &model.CalendarEvent{Title: "Test Event #4 (2019-10-04T12:00:00)", Time: time.Date(2019, 10, 4, 12, 0, 0, 0, time.UTC)}
			err := repo.Create(initialEvent)
			if err != nil {
				t.Fatal(err)
			}

			err = repo.Delete(initialEvent.ID)
			if err != nil {
				t.Fatal(err)
			}

			if _, err = repo.Read(initialEvent.ID); err == nil {
				t.Error("FAILURE: Calendar event is not removed from storage\n")
			}

			repo.Purge()
		})

	t.Run("Checking of EventRepository::ReadAll",
		func(t *testing.T) {

			initialEvents := []*model.CalendarEvent{
				{Title: "Test Event #1 (2019-10-01T12:00:00)", Time: time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC)},
				{Title: "Test Event #2 (2019-10-02T12:00:00)", Time: time.Date(2019, 10, 2, 12, 0, 0, 0, time.UTC)},
				{Title: "Test Event #3 (2019-10-03T12:00:00)", Time: time.Date(2019, 10, 3, 12, 0, 0, 0, time.UTC)},
			}

			for _, event := range initialEvents {
				err := repo.Create(event)
				if err != nil {
					t.Fatal(err)
				}
			}

			resultEvents := repo.ReadAll()
			length := len(resultEvents)
			if length != 3 {
				t.Error("FAILURE: Storage must contain exactly three events\n")
			}

			for i := 0; i < length; i++ {
				if *resultEvents[i] != *initialEvents[i] {
					t.Logf("initial[%d]: %v\n", i, *initialEvents[i])
					t.Logf("result[%d]: %v\n", i, *resultEvents[i])
					t.Error("FAILURE: Initail calendar event is not equals repository event\n")
				}
			}

			repo.Purge()
		})
}
