package repository

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
)

//GroupRepository Storage interface for Banner Group
type GroupRepository interface {
	Create(*model.Group) error
	Read(int64) (*model.Group, error)
	Update(*model.Group) error
	Delete(int64) error
	IsExists(int64) (bool, error)
	GetByCaption(string) (*model.Group, error)
}
