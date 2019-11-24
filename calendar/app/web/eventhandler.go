package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
	"github.com/gzavodov/otus-go/calendar/app/domain/repository"
	"go.uber.org/zap"
)

//EventError event action result error
type EventError struct {
	Error string `json:"error"`
}

//EventIDResult event action result with identifier
type EventIDResult struct {
	Result uint32 `json:"result"`
}

//EventResult event action result with model
type EventResult struct {
	Result *model.Event `json:"result"`
}

//EventListResult event action result with list of models
type EventListResult struct {
	Result []*model.Event `json:"result"`
}

//EventHandler base structure for event action handlers
type EventHandler struct {
	Name   string
	Repo   repository.EventRepository
	Logger *zap.Logger
}

//GetQualifiedName return handler qualified name (used in logging proccess)
func (h *EventHandler) GetQualifiedName() string {
	return fmt.Sprintf("Web::%sHandler", h.Name)
}

//LogRequestURL writes request URL in log
func (h *EventHandler) LogRequestURL(r *http.Request) {
	h.Logger.Info(h.GetQualifiedName(), zap.String("URL Path", r.URL.Path))
}

//LogError writes error in log
func (h *EventHandler) LogError(name string, err error) {
	h.Logger.Error(
		h.GetQualifiedName(),
		zap.NamedError(name, err),
	)
}

//WriteEventResult writes action result in response
func (h *EventHandler) WriteEventResult(w http.ResponseWriter, eventResult interface{}) error {
	output, err := json.Marshal(eventResult)
	if err != nil {
		return fmt.Errorf("could not serialize event result (%w)", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return nil
}

//MethodNotAllowedError writes 405 (StatusMethodNotAllowed) code in response
func (h *EventHandler) MethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	h.Logger.Error(
		h.GetQualifiedName(),
		zap.NamedError("Bad request method", fmt.Errorf("method %s is not allowed", r.Method)),
	)
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
