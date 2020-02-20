package algorithm

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestUCB1Coverage(t *testing.T) {
	count := 1000
	banditArms := make([]BanditArm, 0, count)
	for i := 0; i < count; i++ {
		banditArms = append(banditArms, &BaseBanditArm{})
	}

	alg, err := NewUCB1(banditArms)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < count; i++ {
		index := alg.ResolveArmIndex()
		if index < 0 {
			t.Error(err)
		}

		banditArms[index].SetCount(banditArms[index].GetCount() + 1)
	}

	omittedCount := 0
	for i := 0; i < count; i++ {
		if banditArms[i].GetCount() == 0 {
			omittedCount++
		}
	}

	if omittedCount > 0 {
		t.Fatal(fmt.Errorf("it is expected that all items will be touched after first pass, but %d items are untouched", omittedCount))
	}
}
func TestUCB1Optimality(t *testing.T) {
	armCount := 1000
	banditArms := make([]BanditArm, 0, armCount)
	for i := 0; i < armCount; i++ {
		banditArms = append(banditArms, &BaseBanditArm{})
	}

	alg, err := NewUCB1(banditArms)
	if err != nil {
		t.Fatal(err)
	}

	tryCount := armCount * 10
	for i := 0; i < tryCount; i++ {
		index := alg.ResolveArmIndex()
		if index < 0 {
			t.Error(err)
		}

		banditArms[index].SetCount(banditArms[index].GetCount() + 1)
		rand.Seed(time.Now().UnixNano())
		if rand.Int31n(10) >= 5 {
			banditArms[index].SetReward(banditArms[index].GetReward() + 1)
		}
	}

	var totalCount int64

	var maxCount int64
	var maxCountIndex int

	var maxReward float64
	var maxRewardIndex int

	for i := 0; i < armCount; i++ {
		count := banditArms[i].GetCount()
		if maxCount < count {
			maxCount = count
			maxCountIndex = i
		}

		reward := banditArms[i].GetReward()
		if maxReward < reward {
			maxReward = reward
			maxRewardIndex = i
		}

		totalCount += count
	}

	if totalCount != int64(tryCount) {
		t.Errorf("it is expected that overall count will be equals to total tries count, but overall count is %d and total tries count is:%d\n", totalCount, tryCount)
	}

	if maxCountIndex != maxRewardIndex {
		t.Errorf("it is expected that item with maximum reward will be chosen maximum times, but index of item with maximum reward is %d and index of item with maximum count is:%d\n", maxRewardIndex, maxCountIndex)
	}
}
