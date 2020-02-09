package usecase

import (
	"sync"

	"github.com/gzavodov/otus-go/banner-rotation/algorithm"
	"github.com/gzavodov/otus-go/banner-rotation/model"
	"github.com/gzavodov/otus-go/banner-rotation/repository"
)

func NewBannerUsecase(
	bannerRepo repository.BannerRepository,
	bindingRepo repository.BindingRepository,
	statisticsRepo repository.StatisticsRepository,
	algorithmTypeID int,
) *Banner {
	return &Banner{
		bannerRepo:      bannerRepo,
		bindingRepo:     bindingRepo,
		statisticsRepo:  statisticsRepo,
		algorithmTypeID: algorithmTypeID,
		mu:              sync.RWMutex{},
	}
}

type Banner struct {
	bannerRepo      repository.BannerRepository
	bindingRepo     repository.BindingRepository
	statisticsRepo  repository.StatisticsRepository
	algorithmTypeID int

	mu sync.RWMutex
}

func (c *Banner) Create(m *model.Banner) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.bannerRepo.Create(m)
}

func (c *Banner) Read(ID int64) (*model.Banner, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.bannerRepo.Read(ID)
}

func (c *Banner) Update(m *model.Banner) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.bannerRepo.Update(m)
}

func (c *Banner) Delete(ID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.bannerRepo.Delete(ID)
}

func (c *Banner) AddToSlot(bannerID int64, slotID int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//Check if binding already exists
	binding, err := c.bindingRepo.GetBinding(bannerID, slotID)
	if err != nil {
		return 0, err
	}

	if binding != nil {
		return binding.ID, nil
	}

	binding = &model.Binding{BannerID: bannerID, SlotID: slotID}
	if err := c.bindingRepo.Create(binding); err != nil {
		return 0, err
	}

	return binding.ID, nil
}

func (c *Banner) DeleteFromSlot(bannerID int64, slotID int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//Check if binding exists
	binding, err := c.bindingRepo.GetBinding(bannerID, slotID)
	if err != nil {
		return 0, err
	}

	if binding == nil {
		return 0, nil
	}

	if err := c.bindingRepo.Delete(binding.ID); err != nil {
		return 0, err
	}

	return binding.ID, nil
}

func (c *Banner) RegisterClick(bannerID int64, groupID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.statisticsRepo.IncrementNumberOfClicks(bannerID, groupID)
}

func (c *Banner) Choose(slotID int64, groupID int64) (int64, error) {
	c.mu.RLock()
	statisticsList, err := c.statisticsRepo.GetRotationStatistics(slotID, groupID)
	c.mu.RUnlock()

	if err != nil {
		return -1, err
	}

	quantity := len(statisticsList)
	adapters := make([]algorithm.Statistics, quantity)
	for i := 0; i < quantity; i++ {
		adapters[i] = NewStatisticsAdapter(statisticsList[i])
	}

	alg, err := algorithm.CreateMultiArmedBandit(c.algorithmTypeID, adapters)
	if err != nil {
		return -1, err
	}

	index := alg.ResolveArmIndex()
	if index < 0 {
		return 0, nil
	}

	bannerID := statisticsList[index].BannerID

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.statisticsRepo.IncrementNumberOfShows(bannerID, groupID); err != nil {
		return -1, err
	}
	return bannerID, nil
}
