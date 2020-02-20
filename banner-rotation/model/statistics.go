package model

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

//GetReward implements of algorithm.Arm::GetReward(), returns reward value for multi-armed bandit arm
func (s *Statistics) GetReward() float64 {
	return float64(s.NumberOfClicks)
}

//SetReward implements of algorithm.Arm::SetReward(), assigns reward value for multi-armed bandit arm
func (s *Statistics) SetReward(value float64) {
	s.NumberOfClicks = int64(value)
}
