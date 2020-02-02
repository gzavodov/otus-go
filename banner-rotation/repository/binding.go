package repository

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
)

//BindingRepository Storage interface for Banner Binding
type BindingRepository interface {
	Create(*model.Binding) error
	Read(int64) (*model.Binding, error)
	Update(*model.Binding) error

	Delete(int64) error
	DeleteByModel(m *model.Binding) error
	DeleteByBannerID(bannerID int64) error
	DeleteBySlotID(slotID int64) error

	IsExists(int64) (bool, error)

	GetBinding(bannerID int64, slotID int64) (*model.Binding, error)
	GetBannerBindings(bannerID int64) ([]*model.Binding, error)
	GetSlotBindings(slotID int64) ([]*model.Binding, error)
}
