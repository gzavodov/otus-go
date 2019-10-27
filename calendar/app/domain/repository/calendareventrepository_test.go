package repository

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
)

//Error Constructors
func CouldNotCreateObjectError() error {
	return fmt.Errorf("could not create object in repository")
}

func CouldNotRemoveObjectError() error {
	return errors.New("could not remove object from repository")
}

func ObjectNotMatchedError(expected, received *model.CalendarEvent) error {
	return fmt.Errorf("object before saving in repository is not equal to object after reading from repository; expected: %v, received: %v", *expected, *received)
}

func ObjectListNotMatchedError(expected, received []*model.CalendarEvent) error {
	return fmt.Errorf("quantity of objects before saving in repository is not equal to quantity of objects after reading from repository; expected: %d, received: %d", len(expected), len(received))
}

func TestInMemoryRepository(t *testing.T) {
	repo := NewInMemoryCalendarEventRepository()

	t.Run("CalendarEventRepository::Create",
		func(t *testing.T) {
			source := &model.CalendarEvent{
				Title: "Test Event #1 (2019-10-01T12:00:00)",
				Time:  time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC),
			}
			err := repo.Create(source)
			if err != nil {
				t.Fatal(err)
			}

			if !repo.IsExists(source.ID) {
				t.Error(CouldNotCreateObjectError())
			}

			repo.Purge()
		})

	t.Run("CalendarEventRepository::Read",
		func(t *testing.T) {
			source := &model.CalendarEvent{
				Title: "Test Event #2 (2019-10-02T12:00:00)",
				Time:  time.Date(2019, 10, 2, 12, 0, 0, 0, time.UTC),
			}

			err := repo.Create(source)
			if err != nil {
				t.Fatal(err)
			}

			result, err := repo.Read(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			if *source != *result {
				t.Error(ObjectNotMatchedError(source, result))
			}

			repo.Purge()
		})

	t.Run("CalendarEventRepository::Update",
		func(t *testing.T) {
			source := &model.CalendarEvent{
				Title: "Test Event #3 (2019-10-03T12:00:00)",
				Time:  time.Date(2019, 10, 3, 12, 0, 0, 0, time.UTC),
			}

			err := repo.Create(source)
			if err != nil {
				t.Fatal(err)
			}

			source.Description = "Test Description"
			err = repo.Update(source)
			if err != nil {
				t.Fatal(err)
			}

			result, err := repo.Read(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			if *source != *result {
				t.Error(ObjectNotMatchedError(source, result))
			}

			repo.Purge()
		})

	t.Run("CalendarEventRepository::Delete",
		func(t *testing.T) {
			source := &model.CalendarEvent{
				Title: "Test Event #4 (2019-10-04T12:00:00)",
				Time:  time.Date(2019, 10, 4, 12, 0, 0, 0, time.UTC),
			}
			err := repo.Create(source)
			if err != nil {
				t.Fatal(err)
			}

			err = repo.Delete(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			if repo.IsExists(source.ID) {
				t.Error(CouldNotRemoveObjectError())
			}

			repo.Purge()
		})

	t.Run("CalendarEventRepository::ReadAll",
		func(t *testing.T) {
			sources := []*model.CalendarEvent{
				{
					Title: "Test Event #1 (2019-10-01T12:00:00)",
					Time:  time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC),
				},
				{
					Title: "Test Event #2 (2019-10-02T12:00:00)",
					Time:  time.Date(2019, 10, 2, 12, 0, 0, 0, time.UTC),
				},
				{
					Title: "Test Event #3 (2019-10-03T12:00:00)",
					Time:  time.Date(2019, 10, 3, 12, 0, 0, 0, time.UTC),
				},
			}

			for _, source := range sources {
				err := repo.Create(source)
				if err != nil {
					t.Fatal(err)
				}
			}

			results := repo.ReadAll()

			if len(results) != len(sources) {
				t.Error(ObjectListNotMatchedError(sources, results))
			}

			for i := 0; i < len(results); i++ {
				if *sources[i] != *results[i] {
					t.Error(ObjectNotMatchedError(sources[i], results[i]))
				}
			}

			repo.Purge()
		})
}
