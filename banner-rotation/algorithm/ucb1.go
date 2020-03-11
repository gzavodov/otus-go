package algorithm

import (
	"math"
	"sync"
)

//NewUCB1 creates new multi-armed bandit based on UCB1 algorithm
func NewUCB1(arms []BanditArm) (*UCB1, error) {
	return &UCB1{Arms: arms, mu: sync.RWMutex{}}, nil
}

//UCB1 represents implementation of the upper confidence bound algorithm
type UCB1 struct {
	Arms []BanditArm
	mu   sync.RWMutex
}

const UCB1Factor = 2.0

//ResolveArmIndex resolves an arm index that exploits if the value is more than the epsilon threshold, and explore if the value is less than epsilon
func (b *UCB1) ResolveArmIndex() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var totalCount int64

	quantity := len(b.Arms)
	for i := 0; i < quantity; i++ {
		arm := b.Arms[i]
		if arm.GetCount() == 0 {
			return i
		}

		totalCount += arm.GetCount()
	}

	maxValue := float64(0.0)
	maxIndex := -1

	numerator := UCB1Factor * math.Log(float64(totalCount))
	for i := 0; i < quantity; i++ {
		denominator := float64(b.Arms[i].GetCount())
		value := b.Arms[i].GetAverageReward() + math.Sqrt(numerator/denominator)
		if maxValue < value {
			maxValue = value
			maxIndex = i
		}
	}

	return maxIndex
}

//RegisterArmTouch will increment arm touch count
func (b *UCB1) RegisterArmTouch(armIndex int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if armIndex < 0 || armIndex >= len(b.Arms) {
		return ErrorArmIndexOutOfRange
	}

	b.Arms[armIndex].SetCount(b.Arms[armIndex].GetCount() + 1)
	return nil
}

//RegisterArmReward will update an arm with specified reward value
func (b *UCB1) RegisterArmReward(armIndex int, reward float64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if armIndex < 0 || armIndex >= len(b.Arms) {
		return ErrorArmIndexOutOfRange
	}

	if reward < 0 {
		return ErrorInvalidReward
	}

	return b.Arms[armIndex].AddReward(reward)
}

//GetCounts returns the counts
func (b *UCB1) GetCounts() []int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	quantity := len(b.Arms)
	counts := make([]int64, quantity)
	for i := 0; i < quantity; i++ {
		counts[i] = b.Arms[i].GetCount()
	}
	return counts
}

//GetRewards returns the rewards
func (b *UCB1) GetRewards() []float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	quantity := len(b.Arms)
	rewards := make([]float64, quantity)
	for i := 0; i < quantity; i++ {
		rewards[i] = b.Arms[i].GetAverageReward()
	}
	return rewards
}

func (b *UCB1) GetArm(index int) BanditArm {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if index >= 0 && index <= len(b.Arms) {
		return b.Arms[index]
	}
	return nil
}

func (b *UCB1) AddArm(arm BanditArm) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Arms = append(b.Arms, arm)
}

func (b *UCB1) RemoveArm(index int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if index >= 0 && index <= len(b.Arms) {
		b.Arms = append(
			b.Arms[:index],
			b.Arms[index+1:]...,
		)
	}
}

func (b *UCB1) GetArmCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.Arms)
}
