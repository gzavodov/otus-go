package rest

import (
	"net/http"
	"time"

	"github.com/gzavodov/otus-go/banner-rotation/endpoint"
	"github.com/gzavodov/otus-go/banner-rotation/queue"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
	"go.uber.org/zap"
)

func NewBannerHandler(ucase *usecase.Banner, serviceName string, notificationChannel queue.NotificationChannel, logger *zap.Logger) *BannerHandler {
	return &BannerHandler{
		EntityHandler: EntityHandler{
			Accessor: &BannerAccessor{ucase: ucase},
			Handler: endpoint.Handler{
				Name:                "Banner",
				ServiceName:         serviceName,
				NotificationChannel: notificationChannel,
				Logger:              logger,
			},
		},
	}
}

type BannerHandler struct {
	EntityHandler
}

func (h *BannerHandler) AddToSlot(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != PostMethod {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	bannerID, err := form.ParseInt64("bannerId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slotID, err := form.ParseInt64("slotId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bindingID, err := h.Accessor.(*BannerAccessor).AddToSlot(bannerID, slotID)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.WriteResult(w, bindingID); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *BannerHandler) DeleteFromSlot(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != PostMethod {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	bannerID, err := form.ParseInt64("bannerId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slotID, err := form.ParseInt64("slotId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bindingID, err := h.Accessor.(*BannerAccessor).DeleteFromSlot(bannerID, slotID)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := endpoint.Result{Result: bindingID}
	if err = h.WriteResult(w, result); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *BannerHandler) IsInSlot(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != PostMethod {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	bannerID, err := form.ParseInt64("bannerId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slotID, err := form.ParseInt64("slotId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.Accessor.(*BannerAccessor).IsInSlot(bannerID, slotID)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.WriteResult(w, result); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *BannerHandler) RegisterClick(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != PostMethod {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	bannerID, err := form.ParseInt64("bannerId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	groupID, err := form.ParseInt64("groupId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Accessor.(*BannerAccessor).RegisterClick(bannerID, groupID)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notification := queue.Notification{
		EventType: queue.EventClick,
		BannerID:  bannerID,
		GroupID:   groupID,
		Time:      time.Now().UTC(),
	}

	err = h.Notify(&notification)
	if err != nil {
		h.LogError("Notification", err)
	}

	result := endpoint.Result{Result: "OK"}
	if err = h.WriteResult(w, result); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *BannerHandler) Choose(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != PostMethod {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	slotID, err := form.ParseInt64("slotId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	groupID, err := form.ParseInt64("groupId", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ID, err := h.Accessor.(*BannerAccessor).Choose(slotID, groupID)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notification := queue.Notification{
		EventType: queue.EventChoice,
		BannerID:  ID,
		GroupID:   groupID,
		Time:      time.Now().UTC(),
	}

	err = h.Notify(&notification)
	if err != nil {
		h.LogError("Notification", err)
	}

	if err = h.WriteResult(w, ID); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
