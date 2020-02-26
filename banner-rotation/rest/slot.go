package rest

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
)

type Slot struct {
	ucase *usecase.Slot
}

func (h *Slot) ParseEntity(form *RequestForm) (interface{}, error) {
	return form.ParseSlot()
}

func (h *Slot) ParseEntityIdentity(form *RequestForm) (interface{}, error) {
	return form.ParseInt64("ID", 0)
}

func (h *Slot) ParseEntityCaption(form *RequestForm) (string, error) {
	return form.ParseString("caption", "")
}

func (h *Slot) CreateEntity(entity interface{}) error {
	return h.ucase.Create(entity.(*model.Slot))
}

func (h *Slot) ReadEntity(identity interface{}) (interface{}, error) {
	return h.ucase.Read(identity.(int64))
}

func (h *Slot) UpdateEntity(identity interface{}, entity interface{}) error {
	m := entity.(*model.Slot)
	m.ID = identity.(int64)
	return h.ucase.Update(m)
}

func (h *Slot) DeleteEntity(identity interface{}) error {
	return h.ucase.Delete(identity.(int64))
}

func (h Slot) GetEntityByCaption(caption string) (interface{}, error) {
	return h.ucase.GetByCaption(caption)
}
