package sql

import (
	"context"

	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
)

//SlotRepository Storage interface for Banner Slot
type SlotRepository struct {
	ReferenceRepository
}

//NewSlotRepository creates new SQL Slot Repository
func NewSlotRepository(ctx context.Context, dataSourceName string) repository.SlotRepository {
	if ctx == nil {
		ctx = context.Background()
	}

	return &SlotRepository{
		ReferenceRepository{
			TableName:      "banner_slot",
			BaseRepository: BaseRepository{ctx: ctx, dataSourceName: dataSourceName},
		},
	}
}

//Create creates new Banner Slot in database
//If succseed ID field will be updated
func (r *SlotRepository) Create(slot *model.Slot) error {
	return r.DoCreate(slot)
}

//Read reads Banner Slot from database by ID
func (r *SlotRepository) Read(id int64) (*model.Slot, error) {
	slot := &model.Slot{}
	if err := r.DoRead(id, slot); err != nil {
		return nil, err
	}

	return slot, nil
}

//Update modifies Banner Slot in database
func (r *SlotRepository) Update(slot *model.Slot) error {
	return r.DoUpdate(slot)
}

//Delete removes Banner Slot from database
func (r *SlotRepository) Delete(id int64) error {
	return r.DoDelete(id)
}

//IsExists check if repository contains banner with specified ID
func (r *SlotRepository) IsExists(id int64) (bool, error) {
	return r.CheckIfExists(id)
}

//GetByCaption returns Banner Slot bt specified caption
func (r *SlotRepository) GetByCaption(caption string) (*model.Slot, error) {
	slot := &model.Slot{}
	if err := r.DoGetByCaption(caption, slot); err != nil {
		return nil, err
	}

	return slot, nil
}
