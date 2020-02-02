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

func (c *Group) Read(ID int64) (*model.Group, error) {
	return c.repo.Read(ID)
}

func (c *Group) Update(m *model.Group) error {
	return c.repo.Update(m)
}

func (c *Group) Delete(ID int64) error {
	return c.repo.Delete(ID)
}
