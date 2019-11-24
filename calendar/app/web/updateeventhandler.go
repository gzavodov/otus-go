package web

import (
	"errors"
	"net/http"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"go.uber.org/zap"
)

//NewUpdateEventHandler creates new instance of UpdateEventHandler
func NewUpdateEventHandler(repo repository.EventRepository, logger *zap.Logger) *UpdateEventHandler {
	return &UpdateEventHandler{EventHandler{Name: "UpdateEvent", Repo: repo, Logger: logger}}
}

//UpdateEventHandler request handler for event modification action
type UpdateEventHandler struct {
	EventHandler
}

//ServeHTTP implementation of HandlerFunc::ServeHTTP
func (h UpdateEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}

	eventID, err := form.ParseUint32("ID", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if eventID <= 0 {
		err = errors.New("The event ID must be defined and be greater then zero")
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := form.ParseEvent()
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event.ID = eventID

	if err := h.Repo.Update(event); err != nil {
		h.LogError("Repository", err)
		if err = h.WriteEventResult(w, EventError{Error: err.Error()}); err != nil {
			h.LogError("Response writing", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	eventResult := EventResult{Result: event}
	if err = h.WriteEventResult(w, &eventResult); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
