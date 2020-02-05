package algorithm

import (
	"fmt"
	"testing"
)

func TestUCB1Coverage(t *testing.T) {
	count := 50
	statisticsList := make([]Statistics, 0, count)
	for i := 0; i < count; i++ {
		statisticsList = append(statisticsList, &BaseStatistics{})
	}

	alg, err := NewUCB1(statisticsList)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < count; i++ {
		index := alg.ResolveArmIndex()
		if index < 0 {
			t.Error(err)
		}

		statisticsList[index].SetCount(statisticsList[index].GetCount() + 1)
	}

	omittedCount := 0
	for i := 0; i < count; i++ {
		if statisticsList[i].GetCount() == 0 {
			omittedCount++
		}
	}

	if omittedCount > 0 {
		t.Fatal(fmt.Errorf("is expected that all items will be touched after first pass, but %d item(s) are untouched", omittedCount))
	}
}
