package model

import "fmt"

//RewardValue reward is one click on shown banner
const RewardValue = 1.0

//Statistics represents statistics
type Statistics struct {
	BannerID       int64 `json:"bannerId"`
	GroupID        int64 `json:"groupId"`
	NumberOfShows  int64 `json:"numberOfShows"`
	NumberOfClicks int64 `json:"NumberOfClicks"`
}

//GetCount implements of algorithm.Arm::GetCount(), returns quantity of interactions with multi-armed bandit arm
func (s *Statistics) GetCount() int64 {
	return s.NumberOfShows
}

//SetCount implements of algorithm.Arm::SetCount(), assigns quantity of interactions with multi-armed bandit arm
func (s *Statistics) SetCount(value int64) {
	s.NumberOfShows = value
}

//GetAverageReward implements of algorithm.Arm::GetAverageReward(), returns average reward for multi-armed bandit arm
func (s *Statistics) GetAverageReward() float64 {
	return float64(s.NumberOfClicks) / float64(s.NumberOfShows)
}

//AddReward implements of algorithm.Arm::AddReward(), adds reward value for multi-armed bandit arm
func (s *Statistics) AddReward(value float64) error {
	if value != RewardValue {
		return fmt.Errorf("reward value %f is not allowed in current context, reward must be equals 1", value)
	}
	s.NumberOfClicks++
	return nil
}
