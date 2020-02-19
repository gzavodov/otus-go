package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gzavodov/otus-go/banner-rotation/queue"
	"go.uber.org/zap"
)

//Error event action result error
type Error struct {
	Error string `json:"error"`
}

//Result event action result
type Result struct {
	Result interface{} `json:"result"`
}

//Handler base struct for end point service action handler
type Handler struct {
	Name                string
	ServiceName         string
	NotificationChannel queue.NotificationChannel
	Logger              *zap.Logger
}

//GetQualifiedName return handler qualified name (used in logging proccess)
func (h *Handler) GetQualifiedName() string {
	return fmt.Sprintf("%s::%s", h.ServiceName, h.Name)
}

//LogRequestURL writes request URL in log
func (h *Handler) LogRequestURL(r *http.Request) {
	if h.Logger != nil {
		h.Logger.Info(h.GetQualifiedName(), zap.String("URL Path", r.URL.Path))
	}
}

//LogError writes error in log
func (h *Handler) LogError(name string, err error) {
	if h.Logger != nil {
		h.Logger.Error(
			h.GetQualifiedName(),
			zap.NamedError(name, err),
		)
	}
}

//MethodNotAllowedError writes 405 (StatusMethodNotAllowed) code in response
func (h *Handler) MethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	if h.Logger != nil {
		h.Logger.Error(
			h.GetQualifiedName(),
			zap.NamedError("Bad request method", fmt.Errorf("method %s is not allowed", r.Method)),
		)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

//WriteResult writes action result in response
func (h *Handler) WriteResult(w http.ResponseWriter, result interface{}) error {
	output, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("could not serialize event result (%w)", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(output); err != nil {
		return fmt.Errorf("could not write event result (%w)", err)
	}
	return nil
}

//IsNotificationEnabled check if notifivation enabled
func (h *Handler) IsNotificationEnabled() bool {
	return h.NotificationChannel != nil
}

//Notify send notification
func (h *Handler) Notify(notification *queue.Notification) error {
	if h.NotificationChannel == nil {
		return nil
	}

	return h.NotificationChannel.Write(notification)
}
