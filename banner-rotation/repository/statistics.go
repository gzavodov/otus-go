package repository

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
)

//StatisticsRepository Storage interface for Banner Statistics
type StatisticsRepository interface {
	Create(*model.Statistics) error
	Read(int64, int64) (*model.Statistics, error)
	Update(*model.Statistics) error

	Delete(int64, int64) error
	DeleteByBannerID(int64) error
	DeleteByGroupID(int64) error

	GetBannerStatistics(int64) ([]*model.Statistics, error)
	GetGroupStatistics(int64) ([]*model.Statistics, error)
	GetRotationStatistics(int64, int64) ([]*model.Statistics, error)

	IncrementNumberOfClicks(int64, int64) error
}
