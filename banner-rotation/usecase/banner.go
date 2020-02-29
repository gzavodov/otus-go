package usecase

import (
	"fmt"
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
		algorithmCache:  make(map[int64]map[int64]algorithm.MultiArmedBandit),
		locks:           make(map[string]chan struct{}),
		algMu:           &sync.RWMutex{},
		lockMu:          &sync.RWMutex{},
	}
}

type Banner struct {
	bannerRepo      repository.BannerRepository
	bindingRepo     repository.BindingRepository
	statisticsRepo  repository.StatisticsRepository
	algorithmTypeID int

	locks          map[string]chan struct{}
	algorithmCache map[int64]map[int64]algorithm.MultiArmedBandit
	lockMu         *sync.RWMutex
	algMu          *sync.RWMutex
}

func (c *Banner) Create(m *model.Banner) error {
	return c.bannerRepo.Create(m)
}

func (c *Banner) Read(id int64) (*model.Banner, error) {
	return c.bannerRepo.Read(id)
}

func (c *Banner) Update(m *model.Banner) error {
	return c.bannerRepo.Update(m)
}

func (c *Banner) Delete(id int64) error {
	if err := c.bannerRepo.Delete(id); err != nil {
		return err
	}

	c.algMu.Lock()
	defer c.algMu.Unlock()

	for groupID := range c.algorithmCache {
		for slotID := range c.algorithmCache[groupID] {
			alg := c.algorithmCache[groupID][slotID]

			count := alg.GetArmCount()
			for i := 0; i < count; i++ {
				if alg.GetArm(i).(*model.Statistics).BannerID == id {
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

	c.algMu.Lock()
	defer c.algMu.Unlock()

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

	c.algMu.Lock()
	defer c.algMu.Unlock()

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

	c.algMu.Lock()
	defer c.algMu.Unlock()

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
	key := fmt.Sprintf("%d:%d", slotID, groupID)
	c.waitForLock(key)
	defer c.releaseLock(key)

	c.algMu.RLock()
	var alg algorithm.MultiArmedBandit
	if _, ok := c.algorithmCache[groupID]; ok {
		alg = c.algorithmCache[groupID][slotID]
	}
	c.algMu.RUnlock()

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

		c.algMu.Lock()
		if _, ok := c.algorithmCache[groupID]; !ok {
			c.algorithmCache[groupID] = make(map[int64]algorithm.MultiArmedBandit)
		}
		c.algorithmCache[groupID][slotID] = newAlg
		c.algMu.Unlock()

		alg = newAlg
	}

	index := alg.ResolveArmIndex()
	if index < 0 {
		return 0, nil
	}

	bannerID := alg.GetArm(index).(*model.Statistics).BannerID

	if err := c.incrementNumberOfShows(bannerID, groupID); err != nil {
		return 0, err
	}

	return bannerID, nil
}

func (c *Banner) waitForLock(key string) {
	c.lockMu.Lock()
	ch := c.locks[key]
	if ch == nil {
		ch = make(chan struct{}, 1)
		c.locks[key] = ch
		ch <- struct{}{}
	}
	c.lockMu.Unlock()

	<-ch
}

func (c *Banner) releaseLock(key string) {
	c.lockMu.RLock()
	ch := c.locks[key]
	c.lockMu.RUnlock()

	if ch != nil {
		ch <- struct{}{}
	}
}

func (c *Banner) incrementNumberOfShows(bannerID int64, groupID int64) error {
	if err := c.statisticsRepo.IncrementNumberOfShows(bannerID, groupID); err != nil {
		return err
	}

	c.algMu.Lock()
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
	c.algMu.Unlock()

	return nil
}
