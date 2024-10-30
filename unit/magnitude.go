package unit

import (
	"math"
)

// for a given integer, returns the number of digits and the number magnitude
func Magnitude(num int64) (int, string) {
	if num < 0 {
		digits, mag := Magnitude(-num)
		return digits, "-" + mag
	}
	if num == 0 {
		return 1, "0"
	}
	if num == 1 {
		return 1, "1s"
	}
	digits := int(math.Ceil(math.Log10(float64(num + 1))))
	multipliers := []string{"", "K", "M", "G", "T", "P", "E"}
	multiplier := multipliers[(digits-1)/3]
	tens := (digits - 1) % 3
	result := "1"
	for i := 0; i < tens; i++ {
		result = result + "0"
	}
	result = result + multiplier + "s"
	return digits, result
}
