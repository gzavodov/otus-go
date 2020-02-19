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
		statisticsCache: make(map[int64]map[int64][]*model.Statistics),
	}
}

type Banner struct {
	bannerRepo      repository.BannerRepository
	bindingRepo     repository.BindingRepository
	statisticsRepo  repository.StatisticsRepository
	algorithmTypeID int

	mu              sync.RWMutex
	statisticsCache map[int64]map[int64][]*model.Statistics
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

	return c.remove(ID)
}

func (c *Banner) remove(ID int64) error {
	err := c.bannerRepo.Delete(ID)
	if err != nil {
		return err
	}

	for groupID := range c.statisticsCache {
		for slotID := range c.statisticsCache[groupID] {
			count := len(c.statisticsCache[groupID][slotID])
			for i := 0; i < count; i++ {
				if c.statisticsCache[groupID][slotID][i].BannerID == ID {
					c.statisticsCache[groupID][slotID] = append(
						c.statisticsCache[groupID][slotID][:i],
						c.statisticsCache[groupID][slotID][i+1:]...,
					)
					break
				}
			}
		}
	}

	return nil
}

func (c *Banner) GetByCaption(caption string) (*model.Banner, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.bannerRepo.GetByCaption(caption)
}

func (c *Banner) AddToSlot(bannerID int64, slotID int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.bindToSlot(bannerID, slotID)
}

func (c *Banner) bindToSlot(bannerID int64, slotID int64) (int64, error) {
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

	for groupID := range c.statisticsCache {
		if _, ok := c.statisticsCache[groupID][slotID]; !ok {
			continue
		}

		item, err := c.statisticsRepo.Read(bannerID, groupID)
		if err != nil && !repository.IsNotFoundError(err) {
			return 0, err
		}

		if item == nil {
			item = &model.Statistics{BannerID: bannerID, GroupID: groupID}
		}

		c.statisticsCache[groupID][slotID] = append(
			c.statisticsCache[groupID][slotID],
			item,
		)
	}

	return binding.ID, nil
}

func (c *Banner) DeleteFromSlot(bannerID int64, slotID int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.unbindFromSlot(bannerID, slotID)
}

func (c *Banner) unbindFromSlot(bannerID int64, slotID int64) (int64, error) {
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

	for groupID := range c.statisticsCache {
		if _, ok := c.statisticsCache[groupID][slotID]; !ok {
			continue
		}

		count := len(c.statisticsCache[groupID][slotID])
		for i := 0; i < count; i++ {
			if c.statisticsCache[groupID][slotID][i].BannerID == bannerID {
				c.statisticsCache[groupID][slotID] = append(
					c.statisticsCache[groupID][slotID][:i],
					c.statisticsCache[groupID][slotID][i+1:]...,
				)
				break
			}
		}
	}

	return binding.ID, nil
}

func (c *Banner) IsInSlot(bannerID int64, slotID int64) (bool, error) {

	binding, err := c.bindingRepo.GetBinding(bannerID, slotID)
	if err != nil {
		return false, err
	}

	if binding != nil {
		return true, nil
	}

	return false, nil
}

func (c *Banner) RegisterClick(bannerID int64, groupID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.incrementNumberOfClicks(bannerID, groupID)
}

func (c *Banner) incrementNumberOfClicks(bannerID int64, groupID int64) error {
	if err := c.statisticsRepo.IncrementNumberOfClicks(bannerID, groupID); err != nil {
		return err
	}

	if _, ok := c.statisticsCache[groupID]; ok {
		for slotID := range c.statisticsCache[groupID] {
			count := len(c.statisticsCache[groupID][slotID])
			for i := 0; i < count; i++ {
				if c.statisticsCache[groupID][slotID][i].BannerID == bannerID {
					c.statisticsCache[groupID][slotID][i].NumberOfClicks++
					break
				}
			}
		}
	}

	return nil
}

func (c *Banner) Choose(slotID int64, groupID int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	statisticsList, err := c.getRotationStatistics(slotID, groupID)
	if err != nil {
		return 0, err
	}

	quantity := len(statisticsList)
	adapters := make([]algorithm.Statistics, quantity)
	for i := 0; i < quantity; i++ {
		adapters[i] = NewStatisticsAdapter(statisticsList[i])
	}

	alg, err := algorithm.CreateMultiArmedBandit(c.algorithmTypeID, adapters)
	if err != nil {
		return 0, err
	}

	index := alg.ResolveArmIndex()
	if index < 0 {
		return 0, nil
	}

	bannerID := statisticsList[index].BannerID
	if err := c.incrementNumberOfShows(bannerID, groupID); err != nil {
		return 0, err
	}
	return bannerID, nil
}

func (c *Banner) getRotationStatistics(slotID int64, groupID int64) ([]*model.Statistics, error) {
	if _, ok := c.statisticsCache[groupID]; !ok {
		c.statisticsCache[groupID] = make(map[int64][]*model.Statistics)
	}

	if _, ok := c.statisticsCache[groupID][slotID]; !ok {
		list, err := c.statisticsRepo.GetRotationStatistics(slotID, groupID)
		if err != nil {
			return nil, err
		}
		c.statisticsCache[groupID][slotID] = list
	}

	return c.statisticsCache[groupID][slotID], nil
}

func (c *Banner) incrementNumberOfShows(bannerID int64, groupID int64) error {
	if err := c.statisticsRepo.IncrementNumberOfShows(bannerID, groupID); err != nil {
		return err
	}

	if _, ok := c.statisticsCache[groupID]; !ok {
		return nil
	}

	if _, ok := c.statisticsCache[groupID]; ok {
		for slotID := range c.statisticsCache[groupID] {
			count := len(c.statisticsCache[groupID][slotID])
			for i := 0; i < count; i++ {
				if c.statisticsCache[groupID][slotID][i].BannerID == bannerID {
					c.statisticsCache[groupID][slotID][i].NumberOfShows++
					break
				}
			}
		}
	}

	return nil
}
