package usecase

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
)

func NewGroupUsecase(repo repository.GroupRepository) *Group {
	return &Group{repo: repo}
}

type Group struct {
	repo repository.GroupRepository
}

func (c *Group) Create(m *model.Group) error {
	return c.repo.Create(m)
}

func (c *Group) Read(id int64) (*model.Group, error) {
	return c.repo.Read(id)
}

func (c *Group) Update(m *model.Group) error {
	return c.repo.Update(m)
}

func (c *Group) Delete(id int64) error {
	return c.repo.Delete(id)
}

func (c *Group) GetByCaption(caption string) (*model.Group, error) {
	return c.repo.GetByCaption(caption)
}
