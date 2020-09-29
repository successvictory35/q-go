package vector_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/vector"
)

func ExampleVector_TensorProduct() {
	v := vector.New(1, 0)
	vv := v.TensorProduct(v)
	fmt.Println(vv)

	// Output:
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
}

func ExampleVector_OuterProduct() {
	v := vector.New(1, 0)
	vv := v.OuterProduct(v)
	fmt.Println(vv)

	// Output:
	// [[(1+0i) (0+0i)] [(0+0i) (0+0i)]]
}

func TestVector(t *testing.T) {
	cases := []struct {
		v0, v1       vector.Vector
		innerProduct complex128
		isOrthogonal bool
		isUnit       bool
		norm         complex128
	}{
		{vector.New(1, 0), vector.New(1, 0), complex(1, 0), false, true, complex(1, 0)},
		{vector.New(0, 1), vector.New(0, 1), complex(1, 0), false, true, complex(1, 0)},
		{vector.New(1, 0), vector.New(0, 1), complex(0, 0), true, true, complex(1, 0)},
	}

	for _, c := range cases {
		if c.v0.InnerProduct(c.v1) != c.innerProduct {
			t.Errorf("%v", c.v0.InnerProduct(c.v1))
		}

		if c.v0.IsOrthogonal(c.v1) != c.isOrthogonal {
			t.Errorf("%v", c.v0.IsOrthogonal(c.v1))
		}

		if c.v0.IsUnit() != c.isUnit {
			t.Errorf("%v", c.v0.IsUnit())
		}

		if c.v0.Norm() != c.norm {
			t.Errorf("%v", c.v0.Norm())
		}
	}
}
