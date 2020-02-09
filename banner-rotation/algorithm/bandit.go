package algorithm

import "errors"

var (
	//ErrorInvalidLength Invalid Length Error
	ErrorInvalidLength = errors.New("counts and rewards must be of equal length")
	//ErrorInvalidArmCount Invalid Arm Count Error
	ErrorInvalidArmCount = errors.New("arm count must be greater than zero")
	//ErrorArmIndexOutOfRange Arm Index OutOfRange Error
	ErrorArmIndexOutOfRange = errors.New("arms index is out of range")
	//ErrorInvalidReward Invalid Reward Error
	ErrorInvalidReward = errors.New("reward must be greater than zero")
)

//MultiArmedBandit represents the multi-armed bandit interface
type MultiArmedBandit interface {
	ResolveArmIndex() int
	RegisterArmReward(armIndex int, reward float64) error
	GetCounts() []int64
	GetRewards() []float64
}
