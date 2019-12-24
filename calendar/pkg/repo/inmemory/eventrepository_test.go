package inmemory

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gzavodov/otus-go/calendar/model"
)

//Error Constructors
func CouldNotCreateObjectError() error {
	return fmt.Errorf("could not create object in repository")
}

func CouldNotRemoveObjectError() error {
	return errors.New("could not remove object from repository")
}

func ObjectNotMatchedError(expected, received *model.Event) error {
	return fmt.Errorf("object before saving in repository is not equal to object after reading from repository; expected: %v, received: %v", *expected, *received)
}

func ObjectListNotMatchedError(expected, received []*model.Event) error {
	return fmt.Errorf("quantity of objects before saving in repository is not equal to quantity of objects after reading from repository; expected: %d, received: %d", len(expected), len(received))
}

func TestInMemoryRepository(t *testing.T) {
	repo := NewEventRepository()

	t.Run("EventRepository::Create",
		func(t *testing.T) {
			source := &model.Event{
				Title:     "Test Event #1 (2019-10-01T12:00:00)",
				UserID:    1,
				StartTime: time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2019, 10, 1, 13, 0, 0, 0, time.UTC),
			}
			err := repo.Create(source)
			if err != nil {
				t.Fatal(err)
			}

			ok, err := repo.IsExists(source.ID)
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Error(CouldNotCreateObjectError())
			}

			repo.purge()
		})

	t.Run("EventRepository::Read",
		func(t *testing.T) {
			source := &model.Event{
				Title:     "Test Event #2 (2019-10-02T12:00:00)",
				UserID:    1,
				StartTime: time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2019, 10, 1, 13, 0, 0, 0, time.UTC),
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

			repo.purge()
		})

	t.Run("EventRepository::Update",
		func(t *testing.T) {
			source := &model.Event{
				Title:     "Test Event #3 (2019-10-03T12:00:00)",
				UserID:    1,
				StartTime: time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2019, 10, 1, 13, 0, 0, 0, time.UTC),
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

			repo.purge()
		})

	t.Run("EventRepository::Delete",
		func(t *testing.T) {
			source := &model.Event{
				Title:     "Test Event #4 (2019-10-04T12:00:00)",
				UserID:    1,
				StartTime: time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2019, 10, 1, 13, 0, 0, 0, time.UTC),
			}
			err := repo.Create(source)
			if err != nil {
				t.Fatal(err)
			}

			err = repo.Delete(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			ok, err := repo.IsExists(source.ID)
			if err != nil {
				t.Fatal(err)
			}

			if ok {
				t.Error(CouldNotRemoveObjectError())
			}

			repo.purge()
		})

	t.Run("EventRepository::ReadAll",
		func(t *testing.T) {
			sources := []*model.Event{
				{
					Title:     "Test Event #1 (2019-10-01T12:00:00)",
					UserID:    1,
					StartTime: time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2019, 10, 1, 13, 0, 0, 0, time.UTC),
				},
				{
					Title:     "Test Event #2 (2019-10-02T13:00:00)",
					UserID:    1,
					StartTime: time.Date(2019, 10, 1, 13, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2019, 10, 1, 14, 0, 0, 0, time.UTC),
				},
				{
					Title:     "Test Event #3 (2019-10-03T14:00:00)",
					UserID:    1,
					StartTime: time.Date(2019, 10, 1, 14, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2019, 10, 1, 15, 0, 0, 0, time.UTC),
				},
			}

			for _, source := range sources {
				err := repo.Create(source)
				if err != nil {
					t.Fatal(err)
				}
			}

			results, err := repo.ReadAll()
			if err != nil {
				t.Fatal(err)
			}

			if len(results) != len(sources) {
				t.Error(ObjectListNotMatchedError(sources, results))
			}

			for i := 0; i < len(results); i++ {
				if *sources[i] != *results[i] {
					t.Error(ObjectNotMatchedError(sources[i], results[i]))
				}
			}

			repo.purge()
		})
}
