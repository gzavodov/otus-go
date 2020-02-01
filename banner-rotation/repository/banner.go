package repository

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
)

//BannerRepository Storage interface for Banner
type BannerRepository interface {
	Create(*model.Banner) error
	Read(int64) (*model.Banner, error)
	Update(*model.Banner) error
	Delete(int64) error
	IsExists(int64) (bool, error)
}
