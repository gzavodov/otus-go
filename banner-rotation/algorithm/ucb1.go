package algorithm

import (
	"math"
	"sync"
)

//NewUCB1 creates new multi-armed bandit based on UCB1 algorithm
func NewUCB1(statisticsList []Statistics) (*UCB1, error) {
	return &UCB1{StatisticsList: statisticsList, mu: sync.RWMutex{}}, nil
}

//UCB1 represents implementation of the upper confidence bound algorithm
type UCB1 struct {
	StatisticsList []Statistics
	mu             sync.RWMutex
}

//Initialize will initialize the counts and rewards with the specified number of arms
func (b *UCB1) Initialize(armCount int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if armCount < 1 {
		return ErrorInvalidArmCount
	}

	b.StatisticsList = make([]Statistics, armCount)
	return nil
}

//ResolveArmIndex resolves an arm index that exploits if the value is more than the epsilon threshold, and explore if the value is less than epsilon
func (b *UCB1) ResolveArmIndex() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	quantity := len(b.StatisticsList)

	counts := make([]int64, quantity)
	rewards := make([]float64, quantity)

	for i := 0; i < quantity; i++ {
		stat := b.StatisticsList[i]
		if stat.GetCount() == 0 {
			return i
		}

		counts[i] = stat.GetCount()
		rewards[i] = stat.GetReward()
	}

	values := make([]float64, quantity)
	numerator := 2.0 * math.Log(float64(sum(counts)))
	for i := 0; i < quantity; i++ {
		denominator := float64(counts[i])
		values[i] = rewards[i] + math.Sqrt(numerator/denominator)
	}

	return maxIndex(values)
}

//RegisterArmReward will update an arm with specified reward value,
//e.g. click = 1, no click = 0
func (b *UCB1) RegisterArmReward(armIndex int, reward float64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if armIndex < 0 || armIndex >= len(b.StatisticsList) {
		return ErrorArmIndexOutOfRange
	}

	if reward < 0 {
		return ErrorInvalidReward
	}

	currentCount := b.StatisticsList[armIndex].GetCount()
	currentReward := b.StatisticsList[armIndex].GetReward()

	b.StatisticsList[armIndex].SetCount(currentCount + 1)
	b.StatisticsList[armIndex].SetReward((currentReward*float64(currentCount) + reward) / float64(currentCount+1))

	return nil
}

//GetCounts returns the counts
func (b *UCB1) GetCounts() []int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	quantity := len(b.StatisticsList)
	counts := make([]int64, quantity)
	for i := 0; i < quantity; i++ {
		counts[i] = b.StatisticsList[i].GetCount()
	}
	return counts
}

//GetRewards returns the rewards
func (b *UCB1) GetRewards() []float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	quantity := len(b.StatisticsList)
	rewards := make([]float64, quantity)
	for i := 0; i < quantity; i++ {
		rewards[i] = b.StatisticsList[i].GetReward()
	}
	return rewards
}
