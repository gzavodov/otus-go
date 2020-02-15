package rest

import (
	"errors"
	"net/http"

	"github.com/gzavodov/otus-go/banner-rotation/endpoint"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
)

type Slot struct {
	ucase *usecase.Slot

	endpoint.Handler
}

//Create creates new slot
func (h *Slot) Create(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != "POST" {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	m, err := form.ParseSlot()
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

//Read reads slot by ID
func (h *Slot) Read(w http.ResponseWriter, r *http.Request) {
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

//Update updates slot
func (h *Slot) Update(w http.ResponseWriter, r *http.Request) {
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

	m, err := form.ParseSlot()
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

//Delete deletes slot by ID
func (h *Slot) Delete(w http.ResponseWriter, r *http.Request) {
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
