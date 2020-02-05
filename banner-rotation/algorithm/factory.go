package algorithm

import (
	"fmt"
)

//Algorithm Type
const (
	TypeUnknown       = 0
	TypeUCB1          = 1
	TypeEpsilonGreedy = 2
)

//CreateMultiArmedBandit creates multi armed bandit algorithm
func CreateMultiArmedBandit(typeID int, statisticsList []Statistics) (MultiArmedBandit, error) {
	switch typeID {
	case TypeUCB1:
		return NewUCB1(statisticsList)
	default:
		return nil, fmt.Errorf("algorithm type %d is not supported in current context", typeID)
	}
}
