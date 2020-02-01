package service

import (
	"github.com/gzavodov/otus-go/banner-rotation/algorithm"
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
)

func NewBannerService(bindingRepo repository.BindingRepository, statisticsRepo repository.StatisticsRepository) *BannerService {
	return &BannerService{bindingRepo: bindingRepo, statisticsRepo: statisticsRepo}
}

type BannerService struct {
	bindingRepo    repository.BindingRepository
	statisticsRepo repository.StatisticsRepository
}

func (s *BannerService) RegisterBinding(bannerID int64, slotID int64) error {
	return s.bindingRepo.Create(&model.Binding{BannerID: bannerID, SlotID: slotID})
}

func (s *BannerService) UnregisterBinding(bannerID int64, slotID int64) error {
	return s.bindingRepo.DeleteByModel(&model.Binding{BannerID: bannerID, SlotID: slotID})
}

func (s *BannerService) RegisterClick(bannerID int64, groupID int64) error {
	return s.statisticsRepo.IncrementNumberOfClicks(bannerID, groupID)
}

func (s *BannerService) ResolveBanner(slotID int64, groupID int64) (int64, error) {
	statisticsList, err := s.statisticsRepo.GetRotationStatistics(slotID, groupID)
	if err != nil {
		return -1, err
	}

	quantity := len(statisticsList)
	adapters := make([]algorithm.Statistics, quantity)
	for i := 0; i < quantity; i++ {
		adapters[i] = NewStatisticsAdapter(statisticsList[i])
	}

	bandit, err := algorithm.NewUCB1(adapters)
	if err != nil {
		return -1, err
	}

	index := bandit.ResolveArmIndex()
	if index >= 0 {
		return statisticsList[index].BannerID, nil
	}

	return 0, nil

}
