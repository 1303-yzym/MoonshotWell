package sugar

import (
	"math/rand/v2"
)

// RandomWithProbability 随机概率
// 0.01 -> 1.
func RandomWithProbability(probability float64) bool {
	randNum := rand.Float64()

	return randNum <= probability
}

// RandInt generate random int between [min, max).
func RandInt(minV, maxV int) int {
	if minV == maxV {
		return minV
	}

	if maxV < minV {
		minV, maxV = maxV, minV
	}

	return rand.IntN(maxV-minV) + minV
}
