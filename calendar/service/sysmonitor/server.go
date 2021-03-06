package sysmonitor

import (
	"context"
	"net/http"

	"github.com/gzavodov/otus-go/calendar/pkg/endpoint"
	"github.com/gzavodov/otus-go/calendar/pkg/monitoring"
	"go.uber.org/zap"
)

//NewServer Creates new Healthcheck server
func NewServer(address string, middleware monitoring.Middleware, logger *zap.Logger) *Server {
	return &Server{
		Server:     endpoint.Server{Name: "System Monitoring", Address: address, Logger: logger},
		middleware: middleware,
	}
}

//Server Simple Healthcheck Server for calendar event API
type Server struct {
	HTTPServer *http.Server
	middleware monitoring.Middleware
	endpoint.Server
}

//Start Start handling of Web requests
func (s *Server) Start() error {
	err := http.ListenAndServe(s.Address, s.middleware.PrepareMetricExportHandler())
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
