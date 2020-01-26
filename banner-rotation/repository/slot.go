package repository

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
)

//SlotRepository Storage interface for Banner Slot
type SlotRepository interface {
	Create(*model.Slot) error
	Read(int64) (*model.Slot, error)
	Update(*model.Slot) error
	Delete(int64) error
}
