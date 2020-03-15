package monitoring

import (
	"context"
	"net/http"

	"github.com/gzavodov/otus-go/calendar/pkg/endpoint"
	"github.com/gzavodov/otus-go/calendar/repository"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

//NewServer Creates new Web server
func NewServer(address string, repo repository.EventRepository, logger *zap.Logger) *Server {
	return &Server{
		Server: endpoint.Server{Name: "Monitoring", Address: address, Repo: repo, Logger: logger},
	}
}

//Server Simple Web Server for calendar event API
type Server struct {
	HTTPServer *http.Server

	endpoint.Server
}

//Start Start handling of Web requests
func (s *Server) Start() error {
	serverMux := http.NewServeMux()
	serverMux.Handle("/metrics", promhttp.Handler())
	s.HTTPServer = &http.Server{Addr: s.Address, Handler: serverMux}

	err := s.HTTPServer.ListenAndServe()
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return err
}

//Stop stop handling of Web requests
func (s *Server) Stop() error {
	if s.HTTPServer != nil {
		return s.HTTPServer.Shutdown(context.Background())
	}

	return nil
}
