package algorithm

import (
	"math"
	"sync"
)

//NewUCB1 creates new multi-armed bandit based on UCB1 algorithm
func NewUCB1(counts []uint32, rewards []float64) (*UCB1, error) {
	if len(counts) != len(rewards) {
		return nil, ErrorInvalidLength
	}

	return &UCB1{Counts: counts, Rewards: rewards, mu: sync.RWMutex{}}, nil
}

//UCB1 represents implementation of the upper confidence bound algorithm
type UCB1 struct {
	Counts  []uint32
	Rewards []float64

	mu sync.RWMutex
}

//Initialize will initialize the counts and rewards with the specified number of arms
func (b *UCB1) Initialize(armCount int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if armCount < 1 {
		return ErrorInvalidArmCount
	}

	b.Counts = make([]uint32, armCount)
	b.Rewards = make([]float64, armCount)

	return nil
}

//ResolveArmIndex resolves an arm index that exploits if the value is more than the epsilon threshold, and explore if the value is less than epsilon
func (b *UCB1) ResolveArmIndex() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	armCount := len(b.Counts)
	for i := 0; i < armCount; i++ {
		if b.Counts[i] == 0 {
			return i
		}
	}

	values := make([]float64, armCount)
	numerator := 2.0 * math.Log(float64(sum(b.Counts)))
	for i := 0; i < armCount; i++ {
		denominator := float64(b.Counts[i])
		values[i] = b.Rewards[i] + math.Sqrt(numerator/denominator)
	}

	return maxIndex(values)
}

//RegisterArmReward will update an arm with specified reward value,
//e.g. click = 1, no click = 0
func (b *UCB1) RegisterArmReward(armIndex int, reward float64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if armIndex < 0 || armIndex >= len(b.Rewards) {
		return ErrorArmIndexOutOfRange
	}

	if reward < 0 {
		return ErrorInvalidReward
	}

	currentCount := float64(b.Counts[armIndex])
	currentReward := b.Rewards[armIndex]

	b.Counts[armIndex]++
	b.Rewards[armIndex] = (currentReward*currentCount + reward) / (currentCount + 1)

	return nil
}

//GetCounts returns the counts
func (b *UCB1) GetCounts() []uint32 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]uint32, len(b.Counts))
	copy(result, b.Counts)
	return result
}

//GetRewards returns the rewards
func (b *UCB1) GetRewards() []float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]float64, len(b.Rewards))
	copy(result, b.Rewards)
	return result
}
