package endpoint

import (
	"go.uber.org/zap"
)

//Server base struct for end point service
type Server struct {
	Name    string
	Address string
	Logger  *zap.Logger
}

//GetServiceName returns service name
func (s *Server) GetServiceName() string {
	return s.Name
}
