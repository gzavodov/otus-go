package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	Name        string
	ServiceName string
	Logger      *zap.Logger
}

//GetQualifiedName return handler qualified name (used in logging proccess)
func (h *Handler) GetQualifiedName() string {
	return fmt.Sprintf("%s::%s", h.ServiceName, h.Name)
}

//LogRequestURL writes request URL in log
func (h *Handler) LogRequestURL(r *http.Request) {
	h.Logger.Info(h.GetQualifiedName(), zap.String("URL Path", r.URL.Path))
}

//LogError writes error in log
func (h *Handler) LogError(name string, err error) {
	h.Logger.Error(
		h.GetQualifiedName(),
		zap.NamedError(name, err),
	)
}

//MethodNotAllowedError writes 405 (StatusMethodNotAllowed) code in response
func (h *Handler) MethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	h.Logger.Error(
		h.GetQualifiedName(),
		zap.NamedError("Bad request method", fmt.Errorf("method %s is not allowed", r.Method)),
	)
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

//WriteResult writes action result in response
func (h *Handler) WriteResult(w http.ResponseWriter, result interface{}) error {
	output, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("could not serialize event result (%w)", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return nil
}
