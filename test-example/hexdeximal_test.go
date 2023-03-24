package test_example

import (
	"fmt"
	"math"
	"testing"
)

type Log struct {
	t *testing.T
}

func (l Log) translateHexadecimal(m int, n int, input string) string {
	// to 10 hex
	ans := 0
	radix := 0
	for idx := len(input) - 1; idx >= 0; idx-- {
		cur := input[idx]
		if '0' <= cur && cur < '9' {
			// ⭐️⭐️⭐️ A + 32 = a
			ans += int(cur-'0') * int(math.Pow(float64(m), float64(radix)))
		} else {
			ans += int(cur-'A'+10) * int(math.Pow(float64(m), float64(radix)))
		}
		radix++
	}
	//l.t.Log(ans, m, n)
	// to n
	radix = 0
	result := ""
	for ans > 0 {
		res := ans % n
		// ⭐️
		if res > 9 {
			//  %v is used to format values using the default format, which may not be appropriate for formatting an integer value as a character.
			result = fmt.Sprintf("%c%s", 'A'+res-10, result) // same as +=
		} else {
			result = fmt.Sprintf("%d%s", res, result) // same as +=
		}
		ans /= n
	}

	return result
}

func TestHexDeximal001(t *testing.T) {
	var l Log = Log{t: t}
	t.Log(l.translateHexadecimal(16, 8, "7B"))
	t.Log(l.translateHexadecimal(8, 16, "173"))
}

func TestTranslateHexadecimal(t *testing.T) {
	testCases := []struct {
		m        int
		n        int
		input    string
		expected string
	}{
		{16, 10, "3E8", "1000"},
		{10, 2, "11", "1011"},
		{15, 35, "ABCDEF", "8GNHV"},
		{35, 2, "1WZ", "110111001"},
	}

	var l Log = Log{t: t}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%d_%s", tc.m, tc.n, tc.input), func(t *testing.T) {
			result := l.translateHexadecimal(tc.m, tc.n, tc.input)
			if result != tc.expected {
				t.Errorf("expected %s, but got %s", tc.expected, result)
			}
		})
	}
}
