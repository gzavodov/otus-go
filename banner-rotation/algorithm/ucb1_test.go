package algorithm

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const countIncrement = 1
const awardIncrement = 1.0
const randomizerLimit = 10
const randomizerThreshold = 5

func createBandit(armCount int) (*UCB1, error) {
	banditArms := make([]BanditArm, 0, armCount)
	for i := 0; i < armCount; i++ {
		banditArms = append(banditArms, &BaseBanditArm{})
	}

	return NewUCB1(banditArms)
}

func play(bandit *UCB1, tryCount int) error {
	for i := 0; i < tryCount; i++ {
		index := bandit.ResolveArmIndex()
		if index < 0 {
			return errors.New("could not resolve arm index")
		}

		bandit.Arms[index].SetCount(bandit.Arms[index].GetCount() + countIncrement)

		rand.Seed(time.Now().UnixNano())
		if rand.Int31n(randomizerLimit) >= randomizerThreshold {
			err := bandit.Arms[index].AddReward(awardIncrement)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TestUCB1Coverage(t *testing.T) {
	armCount := 1000

	bandit, err := createBandit(armCount)
	if err != nil {
		t.Fatal(err)
	}

	err = play(bandit, len(bandit.Arms))
	if err != nil {
		t.Error(err)
	}

	t.Logf("	| Arm	| Count	| Reward	|\n")
	omittedArmCount := 0
	for i := 0; i < armCount; i++ {
		count := bandit.Arms[i].GetCount()
		reward := bandit.Arms[i].GetAverageReward()

		t.Logf("	| %d	| %d	| %f	|\n", i, count, reward)
		if count == 0 {
			omittedArmCount++
		}
	}

	if omittedArmCount > 0 {
		t.Fatal(fmt.Errorf("it is expected that all items will be touched after first pass, but %d items are untouched", omittedArmCount))
	}
}

func TestUCB1Optimality(t *testing.T) {
	armCount := 1000

	bandit, err := createBandit(armCount)
	if err != nil {
		t.Fatal(err)
	}

	tryCount := armCount * 10
	err = play(bandit, tryCount)
	if err != nil {
		t.Error(err)
	}

	var totalCount int64

	var maxCount int64
	var maxCountIndex int

	var maxReward float64
	var maxRewardIndex int

	t.Logf("	| Arm	| Count	| Reward	|\n")
	for i := 0; i < armCount; i++ {
		count := bandit.Arms[i].GetCount()

		if maxCount < count {
			maxCount = count
			maxCountIndex = i
		}

		reward := bandit.Arms[i].GetAverageReward()
		if maxReward < reward {
			maxReward = reward
			maxRewardIndex = i
		}

		t.Logf("	| %d	| %d	| %f	|\n", i, count, reward)
		totalCount += count
	}

	if totalCount != int64(tryCount) {
		t.Errorf("it is expected that overall count will be equals to total tries count, but overall count is %d and total tries count is %d\n", totalCount, tryCount)
	}

	if maxCountIndex != maxRewardIndex {
		t.Errorf("it is expected that item with maximum reward will be chosen maximum times, but index of item with maximum reward is %d and index of item with maximum count is %d\n", maxRewardIndex, maxCountIndex)
	}
}
