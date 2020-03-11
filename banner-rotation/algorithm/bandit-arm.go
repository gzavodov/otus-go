package algorithm

//BanditArm represents the multi-armed bandit arm statistics contract
type BanditArm interface {
	GetCount() int64
	SetCount(int64)

	GetAverageReward() float64
	AddReward(float64) error
}

//BaseBanditArm represents default implementation of the multi-armed bandit arm contract
type BaseBanditArm struct {
	count   int64
	rewards []float64
}

//GetCount returns quantity of interactions with current arm
func (s *BaseBanditArm) GetCount() int64 {
	return s.count
}

//SetCount assigns quantity of interactions with current arm
func (s *BaseBanditArm) SetCount(value int64) {
	s.count = value
}

//GetAverageReward returns value of average reward for current arm
func (s *BaseBanditArm) GetAverageReward() float64 {
	var result float64
	for _, reward := range s.rewards {
		result += reward
	}
	return result / float64(s.count)
}

//AddReward appends reward value for current arm
func (s *BaseBanditArm) AddReward(value float64) error {
	s.rewards = append(s.rewards, value)
	return nil
}
