package utils

import "math"

func CountDigits(num int) int {
	return int(math.Log10(float64(num))) + 1
}
