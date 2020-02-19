package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/gzavodov/otus-go/banner-rotation/endpoint"
	"github.com/gzavodov/otus-go/banner-rotation/queue"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
)

type Banner struct {
	ucase *usecase.Banner

	endpoint.Handler
}

//Create creates new banner
func (h *Banner) Create(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	m, err := form.ParseBanner()

	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.ucase.Create(m); err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.WriteResult(w, m); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

//Read reads banner by ID
func (h *Banner) Read(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	ID, err := form.ParseInt64("ID", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ID <= 0 {
		err = errors.New("The ID must be defined and be greater then zero")
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := h.ucase.Read(ID)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.WriteResult(w, m); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

//Update updates banner
func (h *Banner) Update(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}

	ID, err := form.ParseInt64("ID", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ID <= 0 {
		err = errors.New("The ID must be defined and be greater then zero")
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := form.ParseBanner()
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.ID = ID

	if err := h.ucase.Update(m); err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.WriteResult(w, m); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

//Delete deletes banner by ID
func (h *Banner) Delete(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	ID, err := form.ParseInt64("ID", 0)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ID <= 0 {
		err = errors.New("The ID must be defined and be greater then zero")
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.ucase.Delete(ID); err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Banner) AddToSlot(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
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

	bindingID, err := h.ucase.AddToSlot(bannerID, slotID)
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

func (h *Banner) DeleteFromSlot(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
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

	bindingID, err := h.ucase.DeleteFromSlot(bannerID, slotID)
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

func (h *Banner) IsInSlot(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
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

	result, err := h.ucase.IsInSlot(bannerID, slotID)
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

func (h *Banner) RegisterClick(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
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

	if err := h.ucase.RegisterClick(bannerID, groupID); err != nil {
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

	if err := h.Notify(&notification); err != nil {
		h.LogError("Notification", err)
	}

	result := endpoint.Result{Result: "OK"}
	if err = h.WriteResult(w, result); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Banner) Choose(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
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

	ID, err := h.ucase.Choose(slotID, groupID)
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

	if err := h.Notify(&notification); err != nil {
		h.LogError("Notification", err)
	}

	if err = h.WriteResult(w, ID); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

//GetByCaption returns banner by caption
func (h *Banner) GetByCaption(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	caption, err := form.ParseString("caption", "")
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if caption == "" {
		err = errors.New("The caption must be defined and be not empty")
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := h.ucase.GetByCaption(caption)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.WriteResult(w, m); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
