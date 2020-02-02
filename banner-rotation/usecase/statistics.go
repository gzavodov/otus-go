package usecase

import "github.com/gzavodov/otus-go/banner-rotation/model"

//NewStatisticsAdapter creates new statistics adapter for multi-armed bandit algorithm
func NewStatisticsAdapter(statistics *model.Statistics) *StatisticsAdapter {
	return &StatisticsAdapter{statistics: statistics}
}

//StatisticsAdapter wraps model.Statistics and implements  multi-armed bandit arm statistics contract
type StatisticsAdapter struct {
	statistics *model.Statistics
}

//GetBannerID returnd model banner ID
func (s *StatisticsAdapter) GetBannerID() int64 {
	return s.statistics.BannerID
}

//GetGroupID returnd model group ID
func (s *StatisticsAdapter) GetGroupID() int64 {
	return s.statistics.GroupID
}

//GetCount returns quantity of interactions with multi-armed bandit arm
func (s *StatisticsAdapter) GetCount() int64 {
	return s.statistics.NumberOfShows
}

//SetCount assigns quantity of interactions with multi-armed bandit arm
func (s *StatisticsAdapter) SetCount(value int64) {
	s.statistics.NumberOfShows = value
}

//GetReward returns reward value for multi-armed bandit arm
func (s *StatisticsAdapter) GetReward() float64 {
	return float64(s.statistics.NumberOfClicks)
}

//SetReward assigns reward value for multi-armed bandit arm
func (s *StatisticsAdapter) SetReward(value float64) {
	s.statistics.NumberOfClicks = int64(value)
}
