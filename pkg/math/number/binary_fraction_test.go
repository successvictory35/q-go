package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleBinaryFraction() {
	// 0.101 -> 1/2 + 1/8 = 0.5 + 0.125 = 0.625
	f := number.BinaryFraction([]int{1, 0, 1})
	fmt.Println(f)

	// Output:
	// 0.625
}

func TestBinaryFraction(t *testing.T) {
	cases := []struct {
		binary []int
		float  float64
	}{
		{[]int{0, 0, 0}, 0.0},
		{[]int{1, 0, 0}, 0.5},
		{[]int{0, 1, 0}, 0.25},
		{[]int{1, 1, 0}, 0.75},
		{[]int{0, 0, 1}, 0.125},
		{[]int{1, 0, 1}, 0.625},
		{[]int{0, 1, 1}, 0.375},
		{[]int{1, 1, 1}, 0.875},
		{[]int{0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1}, 0.41650390625},
		{[]int{0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}, 0.166656494140625},
	}

	for _, c := range cases {
		result := number.BinaryFraction(c.binary)
		if result == c.float {
			continue
		}

		t.Errorf("expected=%v, actual=%v", c.float, result)
	}
}
