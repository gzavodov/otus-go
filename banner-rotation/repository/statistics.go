package repository

import (
	"github.com/gzavodov/otus-go/banner-rotation/model"
)

//StatisticsRepository Storage interface for Banner Statistics
type StatisticsRepository interface {
	Create(*model.Statistics) error
	Read(int64) (*model.Statistics, error)
	Update(*model.Statistics) error
	Delete(int64) error
}
