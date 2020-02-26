package rest

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
)

type Group struct {
	ucase *usecase.Group
}

func (h *Group) ParseEntity(form *RequestForm) (interface{}, error) {
	return form.ParseGroup()
}

func (h *Group) ParseEntityIdentity(form *RequestForm) (interface{}, error) {
	return form.ParseInt64("ID", 0)
}

func (h *Group) ParseEntityCaption(form *RequestForm) (string, error) {
	return form.ParseString("caption", "")
}

func (h *Group) CreateEntity(entity interface{}) error {
	return h.ucase.Create(entity.(*model.Group))
}

func (h *Group) ReadEntity(identity interface{}) (interface{}, error) {
	return h.ucase.Read(identity.(int64))
}

func (h *Group) UpdateEntity(identity interface{}, entity interface{}) error {
	m := entity.(*model.Group)
	m.ID = identity.(int64)
	return h.ucase.Update(m)
}

func (h *Group) DeleteEntity(identity interface{}) error {
	return h.ucase.Delete(identity.(int64))
}

func (h Group) GetEntityByCaption(caption string) (interface{}, error) {
	return h.ucase.GetByCaption(caption)
}
