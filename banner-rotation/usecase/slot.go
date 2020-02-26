package usecase

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
)

func NewSlotUsecase(repo repository.SlotRepository) *Slot {
	return &Slot{repo: repo}
}

type Slot struct {
	repo repository.SlotRepository
}

func (c *Slot) Create(m *model.Slot) error {
	return c.repo.Create(m)
}

func (c *Slot) Read(id int64) (*model.Slot, error) {
	return c.repo.Read(id)
}

func (c *Slot) Update(m *model.Slot) error {
	return c.repo.Update(m)
}

func (c *Slot) Delete(id int64) error {
	return c.repo.Delete(id)
}

func (c *Slot) GetByCaption(caption string) (*model.Slot, error) {
	return c.repo.GetByCaption(caption)
}
