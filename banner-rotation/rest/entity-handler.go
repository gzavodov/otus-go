package rest

import (
	"errors"
	"net/http"

	"github.com/gzavodov/otus-go/banner-rotation/endpoint"
)

type EntityHandler struct {
	Accessor EntityAccessor
	endpoint.Handler
}

func (h *EntityHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != PostMethod {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	entity, err := h.Accessor.ParseEntity(&form)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Accessor.CreateEntity(entity)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.WriteResult(w, entity)
	if err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *EntityHandler) Read(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != PostMethod {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	identity, err := h.Accessor.ParseEntityIdentity(&form)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entity, err := h.Accessor.ReadEntity(identity)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.WriteResult(w, entity)
	if err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *EntityHandler) Update(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != PostMethod {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}

	identity, err := h.Accessor.ParseEntityIdentity(&form)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*
		if ID <= 0 {
			err = errors.New("parameter ID must be defined and be greater then zero")
			h.LogError("Request parsing", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/

	entity, err := h.Accessor.ParseEntity(&form)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Accessor.UpdateEntity(identity, entity)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.WriteResult(w, entity); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *EntityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != PostMethod {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}

	identity, err := h.Accessor.ParseEntityIdentity(&form)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*
		if ID <= 0 {
			err = errors.New("parameter ID must be defined and be greater then zero")
			h.LogError("Request parsing", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/

	err = h.Accessor.DeleteEntity(identity)
	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EntityHandler) GetByCaption(w http.ResponseWriter, r *http.Request) {
	h.LogRequestURL(r)

	if r.Method != PostMethod {
		h.MethodNotAllowedError(w, r)
		return
	}

	form := RequestForm{Request: r}
	caption, err := h.Accessor.ParseEntityCaption(&form)
	if err != nil {
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if caption == "" {
		err = errors.New("parameter caption must be defined and be not empty")
		h.LogError("Request parsing", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entity, err := h.Accessor.GetEntityByCaption(caption)

	if err != nil {
		h.LogError("Repository", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.WriteResult(w, entity); err != nil {
		h.LogError("Response writing", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
