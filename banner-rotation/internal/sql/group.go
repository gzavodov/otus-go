package sql

import (
	"context"

	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
)

//GroupRepository Storage interface for Banner Group
type GroupRepository struct {
	ReferenceRepository
}

//NewGroupRepository creates new SQL Group Repository
func NewGroupRepository(ctx context.Context, dataSourceName string) repository.GroupRepository {
	if ctx == nil {
		ctx = context.Background()
	}

	return &GroupRepository{
		ReferenceRepository{
			TableName:      "banner_group",
			BaseRepository: BaseRepository{ctx: ctx, dataSourceName: dataSourceName},
		},
	}
}

//Create creates new Banner Group in database
//If succseed ID field will be updated
func (r *GroupRepository) Create(group *model.Group) error {
	return r.DoCreate(group)
}

//Read reads Banner Group from database by ID
func (r *GroupRepository) Read(id int64) (*model.Group, error) {
	group := &model.Group{}
	if err := r.DoRead(id, group); err != nil {
		return nil, err
	}

	return group, nil
}

//Update modifies Banner Group in database
func (r *GroupRepository) Update(group *model.Group) error {
	return r.DoUpdate(group)
}

//Delete removes Banner Group from database
func (r *GroupRepository) Delete(id int64) error {
	return r.DoDelete(id)
}

//IsExists check if repository contains Banner Group with specified ID
func (r *GroupRepository) IsExists(id int64) (bool, error) {
	return r.CheckIfExists(id)
}

//GetByCaption returns Banner Group bt specified caption
func (r *GroupRepository) GetByCaption(caption string) (*model.Group, error) {
	group := &model.Group{}
	if err := r.DoGetByCaption(caption, group); err != nil {
		return nil, err
	}

	return group, nil
}
