package number

import (
	"fmt"
	"math"
)

func BinaryToFloat64(binary []int) float64 {
	var d float64
	for i, b := range binary {
		if b != 0 && b != 1 {
			panic(fmt.Sprintf("invalid input: %v", binary))
		}

		if b == 0 {
			continue
		}

		d = d + math.Pow(0.5, float64(i+1))
	}

	return d
}

func ContinuedFraction(f float64, eps ...float64) []int {
	e := epsilon(eps...)
	if f < e {
		return []int{0}
	}

	list := make([]int, 0)
	r := f
	for {
		t := math.Trunc(r)
		list = append(list, int(t))

		diff := r - t
		if diff < e {
			break
		}

		r = 1.0 / diff
	}

	return list
}

func Convergent(cf []int) (int, int, float64) {
	if len(cf) == 1 {
		return cf[0], 1, float64(cf[0])
	}

	s, r := 1, cf[len(cf)-1]
	for i := 2; i < len(cf); i++ {
		s, r = r, cf[len(cf)-i]*r+s
	}
	s = s + cf[0]*r

	return s, r, float64(s) / float64(r)
}

func epsilon(eps ...float64) float64 {
	if len(eps) > 0 {
		return eps[0]
	}

	return 1e-3
}
