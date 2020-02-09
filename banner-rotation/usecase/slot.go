package usecase

import (
	"sync"

	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
)

func NewSlotUsecase(repo repository.SlotRepository) *Slot {
	return &Slot{repo: repo, mu: sync.RWMutex{}}
}

type Slot struct {
	repo repository.SlotRepository

	mu sync.RWMutex
}

func (c *Slot) Create(m *model.Slot) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.repo.Create(m)
}

func (c *Slot) Read(ID int64) (*model.Slot, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.repo.Read(ID)
}

func (c *Slot) Update(m *model.Slot) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.repo.Update(m)
}

func (c *Slot) Delete(ID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.repo.Delete(ID)
}
