package repository

import "fmt"

//Repository Type
const (
	RepositoryTypeUnknown  = 0
	RepositoryTypeInMemory = 1
)

//CreateEventRepository creates repository by repository type
func CreateEventRepository(typeID int) (EventRepository, error) {
	switch typeID {
	case RepositoryTypeInMemory:
		return NewInMemoryEventRepository(), nil
	default:
		return nil, fmt.Errorf("repository type %d is not supported in current context", typeID)
	}
}
