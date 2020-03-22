package endpoint

import (
	"github.com/gzavodov/otus-go/calendar/repository"
)

//EventServer base struct for event end point
type EventServer struct {
	Repo repository.EventRepository

	Server
}
