package algorithm

import "math"

func sum(values []int64) int64 {
	var total int64

	for _, v := range values {
		total += v
	}
	return total
}

func maxIndex(values []float64) int {
	maxValue := math.Inf(-1)
	maxIndex := -1

	for i, v := range values {
		if v > maxValue {
			maxValue = v
			maxIndex = i
		}
	}
	return maxIndex
}
