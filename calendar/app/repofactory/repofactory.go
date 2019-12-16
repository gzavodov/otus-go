package repofactory

import (
	"context"
	"fmt"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"github.com/gzavodov/otus-go/calendar/app/inmemory"
	"github.com/gzavodov/otus-go/calendar/app/rpc"
	"github.com/gzavodov/otus-go/calendar/app/sqldb"
)

//Event Repository Type
const (
	TypeUnknown  = 0
	TypeInMemory = 1
	TypeSQL      = 2
	TypeRPC      = 3
)

//CreateEventRepository creates calendar event repository by type
func CreateEventRepository(ctx context.Context, typeID int, dataSourceName string) (repository.EventRepository, error) {
	switch typeID {
	case TypeInMemory:
		return inmemory.NewEventRepository(), nil
	case TypeSQL:
		return sqldb.NewEventRepository(ctx, dataSourceName), nil
	case TypeRPC:
		//Try to treat dataSourceName as RPC server address
		return rpc.NewEventRepository(ctx, dataSourceName), nil
	default:
		return nil, fmt.Errorf("repository type %d is not supported in current context", typeID)
	}
}
