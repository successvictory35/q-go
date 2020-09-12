package gate

import (
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/vector"
)

func TestCModExp2op(t *testing.T) {
	v0 := vector.TensorProduct(vector.TensorProductN(vector.New(1, 0), 6), vector.New(0, 1))
	g0 := CNOT(7, 2, 4).Apply(CNOT(7, 2, 5))

	cases := []struct {
		b, a, j, N int
		c          int
		t          []int
		v          vector.Vector
	}{
		{7, 7, 0, 15, 2, []int{4, 5, 6, 7}, v0.Clone().Apply(g0)},
	}

	for _, c := range cases {
		// returns Controlled a^(2^j) mod M of b qubits (c qubits + t qubits )
		a := CModExp2(c.b, c.a, c.j, c.N, c.c, c.t)
		if !a.IsUnitary() {
			t.Errorf("modexp is not unitary")
		}

		if !v0.Apply(a).Equals(c.v) {
			t.Fail()
		}
	}

}

func TestCModExp2(t *testing.T) {
	g1 := CNOT(7, 3, 5).Apply(CCNOT(7, 1, 5, 3)).Apply(CNOT(7, 3, 5))
	g2 := CNOT(7, 4, 6).Apply(CCNOT(7, 1, 6, 4)).Apply(CNOT(7, 4, 6))
	g3 := I(7)

	cases := []struct {
		b, a, j, N int
		c          int
		t          []int
		m          matrix.Matrix
	}{
		{7, 7, 0, 15, 2, []int{4, 5, 6, 7}, nil},
		{7, 7, 1, 15, 1, []int{4, 5, 6, 7}, g1.Apply(g2)},
		{7, 7, 2, 15, 0, []int{4, 5, 6, 7}, g3},
		{8, 4, 0, 21, 2, []int{4, 5, 6, 7, 8}, nil},
	}

	for _, c := range cases {
		// returns Controlled a^(2^j) mod M of b qubits (c qubits + t qubits )
		a := CModExp2(c.b, c.a, c.j, c.N, c.c, c.t)
		if !a.IsUnitary() {
			t.Errorf("modexp is not unitary")
		}

		if c.m != nil && !a.Equals(c.m) {
			t.Fail()
		}
	}
}

func TestInverseU(t *testing.T) {
	m := U(1.0, 1.1, 1.2, 1.3)

	inv := m.Inverse()
	im := m.Apply(inv)

	mm, nn := im.Dimension()
	for i := 0; i < mm; i++ {
		for j := 0; j < nn; j++ {
			if i == j {
				if cmplx.Abs(im[i][j]-complex(1, 0)) > 1e-13 {
					t.Errorf("[%v][%v] %v\n", i, j, im[i][j])
				}
				continue
			}
			if cmplx.Abs(im[i][j]-complex(0, 0)) > 1e-13 {
				t.Errorf("[%v][%v] %v\n", i, j, im[i][j])
			}
		}
	}
}

func TestCZ(t *testing.T) {
	expected := New(
		[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, -1, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 0, -1},
	)

	actual := CZ(3, 0, 2)

	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestCNOT(t *testing.T) {
	expected := New(
		[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 1, 0, 0},
		[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 0, 1},
		[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
	)

	actual := CNOT(3, 0, 2)

	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestControlledNot(t *testing.T) {
	g0 := I().Add(Z()).TensorProduct(I())
	g1 := I().Sub(Z()).TensorProduct(X())
	CN := g0.Add(g1).Mul(0.5)

	if !ControlledNot(2, []int{0}, 1).Equals(CN) {
		t.Error(CN)
	}
}

func TestToffoli(t *testing.T) {
	g := NewSlice(13)
	g[0] = I(2).TensorProduct(H())
	g[1] = I(1).TensorProduct(CNOT(2, 0, 1))
	g[2] = I(2).TensorProduct(T().Dagger())
	g[3] = CNOT(3, 0, 2)
	g[4] = I(2).TensorProduct(T())
	g[5] = I(1).TensorProduct(CNOT(2, 0, 1))
	g[6] = I(2).TensorProduct(T().Dagger())
	g[7] = CNOT(3, 0, 2)
	g[8] = I(1).TensorProduct(T().Dagger()).TensorProduct(T())
	g[9] = CNOT(2, 0, 1).TensorProduct(H())
	g[10] = I(1).TensorProduct(T().Dagger()).TensorProduct(I())
	g[11] = CNOT(2, 0, 1).TensorProduct(I())
	g[12] = T().TensorProduct(S()).TensorProduct(I())

	expected := I(3)
	for _, gate := range g {
		expected = expected.Apply(gate)
	}

	actual := Toffoli(3, 0, 1, 2)
	if !actual.Equals(expected, 1e-13) {
		t.Error(actual)
	}
}

func TestFredkin(t *testing.T) {
	m := make(matrix.Matrix, 8)
	m[0] = []complex128{1, 0, 0, 0, 0, 0, 0, 0}
	m[1] = []complex128{0, 1, 0, 0, 0, 0, 0, 0}
	m[2] = []complex128{0, 0, 1, 0, 0, 0, 0, 0}
	m[3] = []complex128{0, 0, 0, 1, 0, 0, 0, 0}
	m[4] = []complex128{0, 0, 0, 0, 1, 0, 0, 0}
	m[5] = []complex128{0, 0, 0, 0, 0, 0, 1, 0}
	m[6] = []complex128{0, 0, 0, 0, 0, 1, 0, 0}
	m[7] = []complex128{0, 0, 0, 0, 0, 0, 0, 1}

	if !Fredkin(3, 0, 1, 2).Equals(m, 1e-13) {
		t.Fail()
	}
}

func TestIsHermite(t *testing.T) {
	if !H().IsHermite() {
		t.Error(H())
	}

	if !X().IsHermite() {
		t.Error(X())
	}

	if !Y().IsHermite() {
		t.Error(Y())
	}

	if !Z().IsHermite() {
		t.Error(Z())
	}
}

func TestIsUnitary(t *testing.T) {
	if !H().IsUnitary() {
		t.Error(H())
	}

	if !X().IsUnitary() {
		t.Error(X())
	}

	if !Y().IsUnitary() {
		t.Error(Y())
	}

	if !Z().IsUnitary() {
		t.Error(Z())
	}

	u := U(1, 2, 3, 4)
	if !u.IsUnitary() {
		t.Error(u)
	}
}

func TestTrace(t *testing.T) {
	trA := I().Trace()
	if trA != complex(2, 0) {
		t.Error(trA)
	}

	trH := H().Trace()
	if trH != complex(0, 0) {
		t.Error(trH)
	}
}
