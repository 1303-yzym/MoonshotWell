package bloomFilter

import "math"

func GetFalsePositiveRate(biteCount uint, existedCount uint, hashCount uint) float64 {
	return 1 - math.Pow(1.0-math.Exp(-float64(hashCount))*(float64(existedCount)+0.5)/(float64(biteCount)-1), float64(hashCount))
}

// ExpectationsBitmapSize 获取一个合理的位图大小
// n 是预期的元素数量。
// p 是期望的误判率。
// 输出位图大小
func ExpectationsBitmapSize(n int, p float64) int {
	m := -1 * (float64(n) * math.Log(p) / (math.Log(2) * math.Log(2)))

	return int(m)
}
