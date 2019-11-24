package web

import (
	"net/http"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"go.uber.org/zap"
)

//NewCreateEventHandler creates new instance of CreateEventHandler
func NewCreateEventHandler(repo repository.EventRepository, logger *zap.Logger) *CreateEventHandler {
	return &CreateEventHandler{EventHandler{Name: "CreateEvent", Repo: repo, Logger: logger}}
}

//CreateEventHandler request handler for event creation action
type CreateEventHandler struct {
	EventHandler
}

//ServeHTTP implementation of HandlerFunc::ServeHTTP
func (h CreateEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	event, err := form.ParseEvent()
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Repo.Create(event); err != nil {
		h.LogError("Repository", err)
		if err = h.WriteEventResult(w, EventError{Error: err.Error()}); err != nil {
			h.LogError("Response writing", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	eventResult := EventResult{Result: event}
	if err = h.WriteEventResult(w, eventResult); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
