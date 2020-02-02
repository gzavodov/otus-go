package usecase

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
)

func NewBindingUsecase(repo repository.BindingRepository) *Binding {
	return &Binding{repo: repo}
}

type Binding struct {
	repo repository.BindingRepository
}

func (c *Binding) Create(m *model.Binding) error {
	return c.repo.Create(m)
}

func (c *Binding) Read(ID int64) (*model.Binding, error) {
	return c.repo.Read(ID)
}

func (c *Binding) Update(m *model.Binding) error {
	return c.repo.Update(m)
}

func (c *Binding) Delete(ID int64) error {
	return c.repo.Delete(ID)
}

func (c *Binding) DeleteByModel(m *model.Binding) error {
	return c.repo.DeleteByModel(m)
}
