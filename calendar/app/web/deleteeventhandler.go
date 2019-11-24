package web

import (
	"errors"
	"net/http"

	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"go.uber.org/zap"
)

//NewDeleteEventHandler creates new instance of DeleteEventHandler
func NewDeleteEventHandler(repo repository.EventRepository, logger *zap.Logger) *DeleteEventHandler {
	return &DeleteEventHandler{EventHandler{Name: "DeleteEvent", Repo: repo, Logger: logger}}
}

//DeleteEventHandler request handler for event deletion action
type DeleteEventHandler struct {
	EventHandler
}

//ServeHTTP implementation of HandlerFunc::ServeHTTP
func (h DeleteEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if err := h.Repo.Delete(eventID); err != nil {
		h.LogError("Repository", err)
		if err = h.WriteEventResult(w, EventError{Error: err.Error()}); err != nil {
			h.LogError("Response writing", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	eventResult := EventIDResult{Result: eventID}
	if err = h.WriteEventResult(w, eventResult); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
