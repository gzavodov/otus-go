package algorithm

import (
	"fmt"
)

//Algorithm Type
const (
	TypeUCB1 = 1
	//TypeEpsilonGreedy = 2
)

//CreateMultiArmedBandit creates multi armed bandit algorithm
func CreateMultiArmedBandit(typeID int, arms []BanditArm) (MultiArmedBandit, error) {
	switch typeID {
	case TypeUCB1:
		return NewUCB1(arms)
	default:
		return nil, fmt.Errorf("algorithm type %d is not supported in current context", typeID)
	}
}
