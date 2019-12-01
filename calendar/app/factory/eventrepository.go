package factory

import (
	"context"
	"fmt"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"github.com/gzavodov/otus-go/calendar/app/inmemory"
	"github.com/gzavodov/otus-go/calendar/app/sqldb"
)

//Event Repository Type
const (
	EventRepositoryTypeUnknown  = 0
	EventRepositoryTypeInMemory = 1
	EventRepositoryTypeSQL      = 2
)

//CreateEventRepository creates calendar event repository by type
func CreateEventRepository(ctx context.Context, typeID int, dataSourceName string) (repository.EventRepository, error) {
	switch typeID {
	case EventRepositoryTypeInMemory:
		return inmemory.NewEventRepository(), nil
	case EventRepositoryTypeSQL:
		return sqldb.NewEventRepository(ctx, dataSourceName), nil
	default:
		return nil, fmt.Errorf("repository type %d is not supported in current context", typeID)
	}
}
