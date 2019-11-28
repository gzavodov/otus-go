package factory

import (
	"fmt"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"github.com/gzavodov/otus-go/calendar/app/inmemory"
)

//Event Repository Type
const (
	EventRepositoryTypeUnknown  = 0
	EventRepositoryTypeInMemory = 1
)

//CreateEventRepository creates calendar event repository by type
func CreateEventRepository(typeID int) (repository.EventRepository, error) {
	switch typeID {
	case EventRepositoryTypeInMemory:
		return inmemory.NewEventRepository(), nil
	default:
		return nil, fmt.Errorf("repository type %d is not supported in current context", typeID)
	}
}
