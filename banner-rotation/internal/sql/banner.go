package sql

import (
	"context"

	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
)

//BannerRepository  Storage interface for Banner
type BannerRepository struct {
	ReferenceRepository
}

//NewBannerRepository creates new SQL Banner Repository
func NewBannerRepository(ctx context.Context, dataSourceName string) repository.BannerRepository {
	if ctx == nil {
		ctx = context.Background()
	}

	return &BannerRepository{
		ReferenceRepository{
			TableName:      "banner",
			BaseRepository: BaseRepository{ctx: ctx, dataSourceName: dataSourceName},
		},
	}
}

//Create creates new Banner in database
//If succseed ID field will be updated
func (r *BannerRepository) Create(banner *model.Banner) error {
	return r.DoCreate(banner)
}

//Read reads Banner from database by ID
func (r *BannerRepository) Read(id int64) (*model.Banner, error) {
	banner := &model.Banner{}
	if err := r.DoRead(id, banner); err != nil {
		return nil, err
	}

	return banner, nil
}

//Update modifies Banner in database
func (r *BannerRepository) Update(banner *model.Banner) error {
	return r.DoUpdate(banner)
}

//Delete removes Banner from database
func (r *BannerRepository) Delete(id int64) error {
	return r.DoDelete(id)
}

//IsExists check if repository contains banner with specified ID
func (r *BannerRepository) IsExists(id int64) (bool, error) {
	return r.CheckIfExists(id)
}

//GetByCaption returns Banner bt specified caption
func (r *BannerRepository) GetByCaption(caption string) (*model.Banner, error) {
	banner := &model.Banner{}
	if err := r.DoGetByCaption(caption, banner); err != nil {
		return nil, err
	}

	return banner, nil
}
