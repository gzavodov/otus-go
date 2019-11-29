package web

import (
	"context"
	"net/http"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"github.com/gzavodov/otus-go/calendar/app/endpoint"
	"go.uber.org/zap"
)

//NewServer Creates new Web server
func NewServer(address string, repo repository.EventRepository, logger *zap.Logger) *Server {
	return &Server{
		Server: endpoint.Server{Name: "Web", Address: address, Repo: repo, Logger: logger},
	}
}

//Server Simple Web Server for calendar event API
type Server struct {
	HTTPServer   *http.Server
	EventHandler *EventHandler

	endpoint.Server
}

//Start Start handling of Web requests
func (s *Server) Start() error {
	s.EventHandler = &EventHandler{
		Handler: endpoint.Handler{
			Name:        "Main",
			ServiceName: s.Name,
			Repo:        s.Repo,
			Logger:      s.Logger,
		},
	}

	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/create_event", s.EventHandler.Create)
	serverMux.HandleFunc("/update_event", s.EventHandler.Update)
	serverMux.HandleFunc("/delete_event", s.EventHandler.Delete)
	serverMux.HandleFunc("/events_for_day", s.EventHandler.EventsForDay)
	serverMux.HandleFunc("/events_for_week", s.EventHandler.EventsForWeek)
	serverMux.HandleFunc("/events_for_month", s.EventHandler.EventsForMonth)

	s.HTTPServer = &http.Server{Addr: s.Address, Handler: serverMux}

	err := s.HTTPServer.ListenAndServe()
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return err
}

//Stop stop handling of Web requests
func (s *Server) Stop() {
	if s.HTTPServer != nil {
		s.HTTPServer.Shutdown(context.Background())
	}
}
