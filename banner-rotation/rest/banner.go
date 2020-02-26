package rest

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/usecase"
)

type BannerAccessor struct {
	ucase *usecase.Banner
}

func (h *BannerAccessor) ParseEntity(form *RequestForm) (interface{}, error) {
	return form.ParseBanner()
}

func (h *BannerAccessor) ParseEntityIdentity(form *RequestForm) (interface{}, error) {
	return form.ParseInt64("ID", 0)
}

func (h *BannerAccessor) ParseEntityCaption(form *RequestForm) (string, error) {
	return form.ParseString("caption", "")
}

func (h *BannerAccessor) CreateEntity(entity interface{}) error {
	return h.ucase.Create(entity.(*model.Banner))
}

func (h *BannerAccessor) ReadEntity(identity interface{}) (interface{}, error) {
	return h.ucase.Read(identity.(int64))
}

func (h *BannerAccessor) UpdateEntity(identity interface{}, entity interface{}) error {
	m := entity.(*model.Banner)
	m.ID = identity.(int64)
	return h.ucase.Update(m)
}

func (h *BannerAccessor) DeleteEntity(identity interface{}) error {
	return h.ucase.Delete(identity.(int64))
}

func (h *BannerAccessor) GetEntityByCaption(caption string) (interface{}, error) {
	return h.ucase.GetByCaption(caption)
}

func (h *BannerAccessor) AddToSlot(bannerID int64, slotID int64) (int64, error) {
	return h.ucase.AddToSlot(bannerID, slotID)
}

func (h *BannerAccessor) DeleteFromSlot(bannerID int64, slotID int64) (int64, error) {
	return h.ucase.DeleteFromSlot(bannerID, slotID)
}

func (h *BannerAccessor) IsInSlot(bannerID int64, slotID int64) (bool, error) {
	return h.ucase.IsInSlot(bannerID, slotID)
}

func (h *BannerAccessor) RegisterClick(bannerID int64, groupID int64) error {
	return h.ucase.RegisterClick(bannerID, groupID)
}

func (h *BannerAccessor) Choose(slotID int64, groupID int64) (int64, error) {
	return h.ucase.Choose(slotID, groupID)
}
