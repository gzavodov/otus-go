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
		algorithmCache:  make(map[int64]map[int64]algorithm.MultiArmedBandit),
	}
}

type Banner struct {
	bannerRepo      repository.BannerRepository
	bindingRepo     repository.BindingRepository
	statisticsRepo  repository.StatisticsRepository
	algorithmTypeID int

	mu             sync.RWMutex
	algorithmCache map[int64]map[int64]algorithm.MultiArmedBandit
}

func (c *Banner) Create(m *model.Banner) error {
	return c.bannerRepo.Create(m)
}

func (c *Banner) Read(ID int64) (*model.Banner, error) {
	return c.bannerRepo.Read(ID)
}

func (c *Banner) Update(m *model.Banner) error {
	return c.bannerRepo.Update(m)
}

func (c *Banner) Delete(ID int64) error {
	if err := c.bannerRepo.Delete(ID); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for groupID := range c.algorithmCache {
		for slotID := range c.algorithmCache[groupID] {
			alg := c.algorithmCache[groupID][slotID]

			count := alg.GetArmCount()
			for i := 0; i < count; i++ {
				if alg.GetArm(i).(*model.Statistics).BannerID == ID {
					alg.RemoveArm(i)
					break
				}
			}
		}
	}

	return nil
}

func (c *Banner) GetByCaption(caption string) (*model.Banner, error) {
	return c.bannerRepo.GetByCaption(caption)
}

func (c *Banner) AddToSlot(bannerID int64, slotID int64) (int64, error) {
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

	c.mu.Lock()
	defer c.mu.Unlock()

	for groupID := range c.algorithmCache {
		alg, ok := c.algorithmCache[groupID][slotID]
		if !ok {
			continue
		}

		stat, err := c.statisticsRepo.Read(bannerID, groupID)
		if err != nil && !repository.IsNotFoundError(err) {
			return 0, err
		}

		if stat == nil {
			stat = &model.Statistics{BannerID: bannerID, GroupID: groupID}
		}

		alg.AddArm(stat)
	}

	return binding.ID, nil
}

func (c *Banner) DeleteFromSlot(bannerID int64, slotID int64) (int64, error) {
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

	c.mu.Lock()
	defer c.mu.Unlock()

	for groupID := range c.algorithmCache {
		alg, ok := c.algorithmCache[groupID][slotID]
		if !ok {
			continue
		}

		count := alg.GetArmCount()
		for i := 0; i < count; i++ {
			if alg.GetArm(i).(*model.Statistics).BannerID == bannerID {
				alg.RemoveArm(i)
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
	if err := c.statisticsRepo.IncrementNumberOfClicks(bannerID, groupID); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.algorithmCache[groupID]; !ok {
		return nil
	}

	for slotID := range c.algorithmCache[groupID] {
		alg := c.algorithmCache[groupID][slotID]

		count := alg.GetArmCount()
		for i := 0; i < count; i++ {
			stat := alg.GetArm(i).(*model.Statistics)
			if stat.BannerID == bannerID {
				stat.NumberOfClicks++
				break
			}
		}
	}

	return nil
}

func (c *Banner) Choose(slotID int64, groupID int64) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var alg algorithm.MultiArmedBandit
	if _, ok := c.algorithmCache[groupID]; ok {
		alg = c.algorithmCache[groupID][slotID]
	}

	if alg == nil {
		statList, err := c.statisticsRepo.GetRotationStatistics(slotID, groupID)
		if err != nil {
			return 0, err
		}

		count := len(statList)
		algStatList := make([]algorithm.BanditArm, count)
		for i := 0; i < count; i++ {
			algStatList[i] = statList[i]
		}

		newAlg, err := algorithm.CreateMultiArmedBandit(c.algorithmTypeID, algStatList)
		if err != nil {
			return 0, err
		}

		if _, ok := c.algorithmCache[groupID]; !ok {
			c.algorithmCache[groupID] = make(map[int64]algorithm.MultiArmedBandit)
		}
		c.algorithmCache[groupID][slotID] = newAlg
		alg = newAlg
	}

	index := alg.ResolveArmIndex()
	if index < 0 {
		return 0, nil
	}

	bannerID := alg.GetArm(index).(*model.Statistics).BannerID
	if err := c.statisticsRepo.IncrementNumberOfShows(bannerID, groupID); err != nil {
		return 0, err
	}

	if _, ok := c.algorithmCache[groupID]; ok {
		for slotID := range c.algorithmCache[groupID] {
			alg := c.algorithmCache[groupID][slotID]

			count := alg.GetArmCount()
			for i := 0; i < count; i++ {
				stat := alg.GetArm(i).(*model.Statistics)
				if stat.BannerID == bannerID {
					stat.NumberOfShows++
					break
				}
			}
		}
	}

	return bannerID, nil
}
