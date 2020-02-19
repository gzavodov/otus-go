package usecase

import (
	"sync"

	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
)

func NewGroupUsecase(repo repository.GroupRepository) *Group {
	return &Group{repo: repo, mu: sync.RWMutex{}}
}

type Group struct {
	repo repository.GroupRepository

	mu sync.RWMutex
}

func (c *Group) Create(m *model.Group) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.repo.Create(m)
}

func (c *Group) Read(ID int64) (*model.Group, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.repo.Read(ID)
}

func (c *Group) Update(m *model.Group) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.repo.Update(m)
}

func (c *Group) Delete(ID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.repo.Delete(ID)
}

func (c *Group) GetByCaption(caption string) (*model.Group, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.repo.GetByCaption(caption)
}
