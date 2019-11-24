package web

import (
	"net/http"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"go.uber.org/zap"
)

//NewEventsForMonthHandler creates new instance of EventsForMonthHandler
func NewEventsForMonthHandler(repo repository.EventRepository, logger *zap.Logger) *EventsForMonthHandler {
	return &EventsForMonthHandler{EventHandler{Name: "EventsForMonth", Repo: repo, Logger: logger}}
}

//EventsForMonthHandler request handler for getting list of month events
type EventsForMonthHandler struct {
	EventHandler
}

//ServeHTTP implementation of HandlerFunc::ServeHTTP
func (h EventsForMonthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "GET" {
		h.MethodNotAllowedError(w, r)
		return
	}

	query := RequestQuery{Request: r}
	date, err := query.ParseDate("date", time.Now())
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 0, time.UTC)

	userID, err := query.ParseUint32("userId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := h.Repo.ReadList(userID, from, to)
	if err != nil {
		h.LogError("Repository", err)
		if err = h.WriteEventResult(w, EventError{Error: err.Error()}); err != nil {
			h.LogError("Response Writing", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = h.WriteEventResult(w, EventListResult{Result: list})
	if err != nil {
		h.LogError("Response Writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
