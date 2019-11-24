package web

import (
	"net/http"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"go.uber.org/zap"
)

//NewServer Creates new web server
func NewServer(address string, repo repository.EventRepository, logger *zap.Logger) *Server {
	server := &Server{Address: address, Repo: repo, Logger: logger}
	return server
}

//Server Simple Web Server for calendar event API
type Server struct {
	Address    string
	Repo       repository.EventRepository
	Logger     *zap.Logger
	HTTPServer *http.ServeMux
}

//Start Start handling of web requests
func (h *Server) Start() error {
	h.HTTPServer = http.NewServeMux()

	h.HTTPServer.Handle("/create_event", NewCreateEventHandler(h.Repo, h.Logger))
	h.HTTPServer.Handle("/update_event", NewUpdateEventHandler(h.Repo, h.Logger))
	h.HTTPServer.Handle("/delete_event", NewDeleteEventHandler(h.Repo, h.Logger))
	h.HTTPServer.Handle("/events_for_day", NewEventsForDayHandler(h.Repo, h.Logger))
	h.HTTPServer.Handle("/events_for_week", NewEventsForWeekHandler(h.Repo, h.Logger))
	h.HTTPServer.Handle("/events_for_month", NewEventsForMonthHandler(h.Repo, h.Logger))

	return http.ListenAndServe(h.Address, h.HTTPServer)
}
