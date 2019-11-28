package rpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"github.com/gzavodov/otus-go/calendar/app/endpoint"
	"go.uber.org/zap"
)

//NewServer Creates new GRPC server
func NewServer(address string, repo repository.EventRepository, logger *zap.Logger) *Server {
	return &Server{
		Server: endpoint.Server{Name: "GRPC", Address: address, Repo: repo, Logger: logger},
	}
}

//Server Simple GPRC Server for calendar event API
type Server struct {
	GRPCServer   *grpc.Server
	EventHandler *EventHandler

	endpoint.Server
}

//Start start handling of GRPC requests
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	s.EventHandler = &EventHandler{
		Handler: endpoint.Handler{
			Name:        "Main",
			ServiceName: s.Name,
			Repo:        s.Repo,
			Logger:      s.Logger,
		},
	}
	s.GRPCServer = grpc.NewServer()
	RegisterEventServiceServer(s.GRPCServer, s.EventHandler)

	s.SetIsStarted(true)
	defer s.SetIsStarted(false)

	return s.GRPCServer.Serve(listener)
}

//Stop stop handling of GRPC requests
func (s *Server) Stop() {
	if s.GRPCServer != nil {
		s.GRPCServer.Stop()
	}
}
