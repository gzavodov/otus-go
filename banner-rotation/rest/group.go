package rest

import (
	"errors"
	"net/http"

	"github.com/gzavodov/otus-go/banner-rotation/endpoint"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
)

type Group struct {
	ucase *usecase.Group

	endpoint.Handler
}

//Create creates new group
func (h *Group) Create(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	m, err := form.ParseGroup()
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

//Read reads group by ID
func (h *Group) Read(w http.ResponseWriter, r *http.Request) {
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

//Update updates group
func (h *Group) Update(w http.ResponseWriter, r *http.Request) {
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

	m, err := form.ParseGroup()
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

//Delete deletes group by ID
func (h *Group) Delete(w http.ResponseWriter, r *http.Request) {
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

//GetByCaption returns group by caption
func (h *Group) GetByCaption(w http.ResponseWriter, r *http.Request) {
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
