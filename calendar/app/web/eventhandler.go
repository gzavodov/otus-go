package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gzavodov/otus-go/calendar/app/domain/model"
	"github.com/gzavodov/otus-go/calendar/app/endpoint"
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
	endpoint.Handler
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

//Create creates new calendar event
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	if err := h.ValidateEvent(event); err != nil {
		h.LogError("Validation", err)
		if err = h.WriteEventResult(w, EventError{Error: err.Error()}); err != nil {
			h.LogError("Response writing", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
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

//Update updates calendar event
func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	if err := h.ValidateEvent(event); err != nil {
		h.LogError("Validation", err)
		if err = h.WriteEventResult(w, EventError{Error: err.Error()}); err != nil {
			h.LogError("Response writing", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

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

//Delete deletes calendar event by ID
func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

//EventsForDay returns calendar events for specified date
func (h *EventHandler) EventsForDay(w http.ResponseWriter, r *http.Request) {
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

	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	to := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC)

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

//EventsForWeek returns calendar events for week of specified date
func (h *EventHandler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
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

	//week starts from monday
	dayIndex := int(date.Weekday())
	if dayIndex > 0 {
		dayIndex--
	} else {
		dayIndex = 6
	}
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1*dayIndex)
	to := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC).AddDate(0, 0, 6-dayIndex)

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

//EventsForMonth returns calendar events for month of specified date
func (h *EventHandler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
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

//ValidateEvent validates calendar event model
func (h *EventHandler) ValidateEvent(event *model.Event) error {
	validator := model.EventValidator{Event: event}
	return validator.Validate().Error()
}
