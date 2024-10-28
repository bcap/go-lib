package unit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMagnitude(t *testing.T) {
	test := func(num float64, expectedDigits int, expectedMag string) {
		digits, mag := Magnitude(num)
		require.Equal(t, expectedDigits, digits)
		require.Equal(t, expectedMag, mag)
		if num != 0 {
			negDigits, negMag := Magnitude(-num)
			require.Equal(t, expectedDigits, negDigits)
			require.Equal(t, "-"+expectedMag, negMag)
		}
	}

	test(0, 1, "0")
	test(0.1, 1, "1s")
	test(0.00001, 1, "1s")
	test(1, 1, "1s")
	test(2, 1, "1s")
	test(9, 1, "1s")
	test(10, 2, "10s")
	test(11, 2, "10s")
	test(99, 2, "10s")
	test(100, 3, "100s")
	test(999, 3, "100s")
	test(1000, 4, "1Ks")
	test(1001, 4, "1Ks")
	test(10001, 5, "10Ks")
	test(1000*1000-1, 6, "100Ks")
	test(1000*1000, 7, "1Ms")
}
