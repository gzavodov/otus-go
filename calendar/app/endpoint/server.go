package endpoint

import (
	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"go.uber.org/zap"
)

//Server base struct for end point services
type Server struct {
	Name      string
	Address   string
	Repo      repository.EventRepository
	Logger    *zap.Logger
	isStarted bool
}

//GetServiceName returns service name
func (s *Server) GetServiceName() string {
	return s.Name
}

//SetIsStarted set server IsStarted flag
func (s *Server) SetIsStarted(isStarted bool) {
	s.isStarted = isStarted
}

//IsStarted check if service started
func (s *Server) IsStarted() bool {
	return s.isStarted
}
