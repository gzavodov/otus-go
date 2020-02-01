package algorithm

//Statistics represents the multi-armed bandit arm statistics contract
type Statistics interface {
	GetCount() int64
	SetCount(int64)

	GetReward() float64
	SetReward(float64)
}

//BaseStatistics represents default implementation of the multi-armed bandit arm statistics contract
type BaseStatistics struct {
	count  int64
	reward float64
}

//GetCount returns quantity of interactions with current arm
func (s *BaseStatistics) GetCount() int64 {
	return s.count
}

//SetCount assigns quantity of interactions with current arm
func (s *BaseStatistics) SetCount(value int64) {
	s.count = value
}

//GetReward returns reward value for current arm
func (s *BaseStatistics) GetReward() float64 {
	return s.reward
}

//SetReward assigns reward value for current arm
func (s *BaseStatistics) SetReward(value float64) {
	s.reward = value
}
