package endpoint

import (
	"fmt"

	"github.com/gzavodov/otus-go/calendar/repository"
	"go.uber.org/zap"
)

//Handler base struct for end point service action handler
type Handler struct {
	Name        string
	ServiceName string
	Repo        repository.EventRepository
	Logger      *zap.Logger
}

//GetQualifiedName return handler qualified name (used in logging proccess)
func (h *Handler) GetQualifiedName() string {
	return fmt.Sprintf("%s::%s", h.ServiceName, h.Name)
}

//LogError writes error in log
func (h *Handler) LogError(name string, err error) {
	h.Logger.Error(
		h.GetQualifiedName(),
		zap.NamedError(name, err),
	)
}
