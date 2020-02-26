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

	quantity := len(b.Arms)

	counts := make([]int64, quantity)
	rewards := make([]float64, quantity)

	for i := 0; i < quantity; i++ {
		arm := b.Arms[i]
		if arm.GetCount() == 0 {
			return i
		}

		counts[i] = arm.GetCount()
		rewards[i] = arm.GetReward()
	}

	values := make([]float64, quantity)
	numerator := UCB1Factor * math.Log(float64(sum(counts)))
	for i := 0; i < quantity; i++ {
		denominator := float64(counts[i])
		values[i] = rewards[i] + math.Sqrt(numerator/denominator)
	}

	return maxIndex(values)
}

//RegisterArmTouch will increment arm touch count
func (b *UCB1) RegisterArmTouch(armIndex int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if armIndex < 0 || armIndex >= len(b.Arms) {
		return ErrorArmIndexOutOfRange
	}

	currentReward := b.Arms[armIndex].GetReward()
	currentCount := b.Arms[armIndex].GetCount()

	b.Arms[armIndex].SetReward(currentReward * float64(currentCount) / float64(currentCount+1))
	b.Arms[armIndex].SetCount(currentCount + 1)
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

	currentReward := b.Arms[armIndex].GetReward()
	currentCount := b.Arms[armIndex].GetCount()

	b.Arms[armIndex].SetReward((currentReward*float64(currentCount) + reward) / float64(currentCount))

	return nil
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
		rewards[i] = b.Arms[i].GetReward()
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
