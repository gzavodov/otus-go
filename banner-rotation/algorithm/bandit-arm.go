package algorithm

//BanditArm represents the multi-armed bandit arm statistics contract
type BanditArm interface {
	GetCount() int64
	SetCount(int64)

	GetReward() float64
	SetReward(float64)
}

//BaseBanditArm represents default implementation of the multi-armed bandit arm contract
type BaseBanditArm struct {
	count  int64
	reward float64
}

//GetCount returns quantity of interactions with current arm
func (s *BaseBanditArm) GetCount() int64 {
	return s.count
}

//SetCount assigns quantity of interactions with current arm
func (s *BaseBanditArm) SetCount(value int64) {
	s.count = value
}

//GetReward returns reward value for current arm
func (s *BaseBanditArm) GetReward() float64 {
	return s.reward
}

//SetReward assigns reward value for current arm
func (s *BaseBanditArm) SetReward(value float64) {
	s.reward = value
}
