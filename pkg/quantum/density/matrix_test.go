package density

import (
	"fmt"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func TestPartialTrace(t *testing.T) {
	qc := gate.H().TensorProduct(gate.I()).Apply(gate.CNOT(2, 0, 1))
	q := qubit.Zero(2).Apply(qc)
	rho := New().Add(1.0, q)

	pt := rho.PartialTrace(0)
	fmt.Println(pt)
}

func TestDensityMatrix(t *testing.T) {
	cases := []struct {
		p         []float64
		q         []*qubit.Qubit
		tr, str   complex128
		expectedM matrix.Matrix
		expectedV complex128
	}{
		{
			[]float64{0.1, 0.9},
			[]*qubit.Qubit{qubit.Zero(), qubit.One()},
			complex(1, 0), complex(0.82, 0),
			gate.X(), 0.0,
		},
		{
			[]float64{0.1, 0.9},
			[]*qubit.Qubit{qubit.Zero(), qubit.Zero().Apply(gate.H())},
			complex(1, 0), complex(0.91, 0),
			gate.X(), 0.9,
		},
	}

	for _, c := range cases {
		rho := New()
		for i := range c.p {
			rho.Add(c.p[i], c.q[i])
		}

		if cmplx.Abs(rho.Trace()-c.tr) > 1e-13 {
			t.Errorf("trace=%v", rho.Trace())
		}

		if cmplx.Abs(rho.Squared().Trace()-c.str) > 1e-13 {
			t.Errorf("strace%v", rho.Squared().Trace())
		}

		if cmplx.Abs(rho.ExpectedValue(c.expectedM)-c.expectedV) > 1e-13 {
			t.Errorf("expected value=%v", rho.ExpectedValue(c.expectedM))
		}
	}
}
