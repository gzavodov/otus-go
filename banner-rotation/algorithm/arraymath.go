package algorithm

import "math"

func sum(values []uint32) uint32 {
	var total uint32

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
