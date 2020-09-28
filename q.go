package q

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
	"strings"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/rand"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

type Qubit int

func (q Qubit) Index() int {
	return int(q)
}

func Index(qb ...Qubit) []int {
	index := make([]int, 0)
	for i := range qb {
		index = append(index, qb[i].Index())
	}

	return index
}

type Q struct {
	internal *qubit.Qubit
	Seed     []int64
	Rand     func(seed ...int64) float64
}

func New() *Q {
	return &Q{
		internal: nil,
		Rand:     rand.Crypto,
	}
}

func (q *Q) New(z ...complex128) Qubit {
	if q.internal == nil {
		q.internal = qubit.New(z...)
		q.internal.Seed = q.Seed
		q.internal.Rand = q.Rand
		return Qubit(0)
	}

	q.internal.TensorProduct(qubit.New(z...))
	index := q.NumberOfBit() - 1
	return Qubit(index)
}

func (q *Q) Zero() Qubit {
	return q.New(1, 0)
}

func (q *Q) One() Qubit {
	return q.New(0, 1)
}

func (q *Q) ZeroWith(bit int) []Qubit {
	r := make([]Qubit, 0)
	for i := 0; i < bit; i++ {
		r = append(r, q.Zero())
	}

	return r
}

func (q *Q) OneWith(bit int) []Qubit {
	r := make([]Qubit, 0)
	for i := 0; i < bit; i++ {
		r = append(r, q.One())
	}

	return r
}

func (q *Q) ZeroLog2(N int) []Qubit {
	n := int(math.Log2(float64(N))) + 1
	return q.ZeroWith(n)
}

func (q *Q) NumberOfBit() int {
	return q.internal.NumberOfBit()
}

func (q *Q) Amplitude() []complex128 {
	return q.internal.Amplitude()
}

func (q *Q) Probability() []float64 {
	return q.internal.Probability()
}

func (q *Q) Measure(qb ...Qubit) *qubit.Qubit {
	if len(qb) < 1 {
		m := make([]*qubit.Qubit, 0)
		for i := 0; i < q.NumberOfBit(); i++ {
			m = append(m, q.internal.Measure(i))
		}

		return qubit.TensorProduct(m...)
	}

	m := make([]*qubit.Qubit, 0)
	m = append(m, q.internal.Measure(qb[0].Index()))
	for i := 1; i < len(qb); i++ {
		m = append(m, q.internal.Measure(qb[i].Index()))
	}

	return qubit.TensorProduct(m...)
}

func (q *Q) MeasureAsInt(qb ...Qubit) int {
	b := q.BinaryString(qb...)
	i, err := strconv.ParseInt(b, 2, 0)
	if err != nil {
		panic(err)
	}

	return int(i)
}

func (q *Q) MeasureAsBinary(qb ...Qubit) []int {
	b := make([]int, 0)
	for _, i := range qb {
		if q.Measure(i).IsZero() {
			b = append(b, 0)
			continue
		}

		b = append(b, 1)
	}

	return b
}

func (q *Q) BinaryString(qb ...Qubit) string {
	var sb strings.Builder
	for _, i := range qb {
		if q.Measure(i).IsZero() {
			sb.WriteString("0")
			continue
		}

		sb.WriteString("1")
	}

	return sb.String()
}

func (q *Q) I(qb ...Qubit) *Q {
	return q.Apply(gate.I(), qb...)
}

func (q *Q) H(qb ...Qubit) *Q {
	return q.Apply(gate.H(), qb...)
}

func (q *Q) X(qb ...Qubit) *Q {
	return q.Apply(gate.X(), qb...)
}

func (q *Q) Y(qb ...Qubit) *Q {
	return q.Apply(gate.Y(), qb...)
}

func (q *Q) Z(qb ...Qubit) *Q {
	return q.Apply(gate.Z(), qb...)
}

func (q *Q) S(qb ...Qubit) *Q {
	return q.Apply(gate.S(), qb...)
}

func (q *Q) T(qb ...Qubit) *Q {
	return q.Apply(gate.T(), qb...)
}

func (q *Q) U(alpha, beta, gamma, delta float64, qb ...Qubit) *Q {
	return q.Apply(gate.U(alpha, beta, gamma, delta), qb...)
}

func (q *Q) RX(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.RX(theta), qb...)
}

func (q *Q) RY(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.RY(theta), qb...)
}

func (q *Q) RZ(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.RZ(theta), qb...)
}

func (q *Q) Apply(mat matrix.Matrix, qb ...Qubit) *Q {
	if len(qb) < 1 {
		q.internal.Apply(mat)
		return q
	}

	index := Index(qb...)

	g := gate.I()
	if index[0] == 0 {
		g = mat
	}

	for i := 1; i < q.NumberOfBit(); i++ {
		found := false
		for j := range index {
			if i == index[j] {
				found = true
				break
			}
		}

		if found {
			g = g.TensorProduct(mat)
			continue
		}

		g = g.TensorProduct(gate.I())
	}

	q.internal.Apply(g)
	return q
}

func (q *Q) ControlledR(control []Qubit, target Qubit, k int) *Q {
	cr := gate.ControlledR(q.NumberOfBit(), Index(control...), target.Index(), k)
	q.internal.Apply(cr)
	return q
}

func (q *Q) CR(control, target Qubit, k int) *Q {
	return q.ControlledR([]Qubit{control}, target, k)
}

func (q *Q) ControlledZ(control []Qubit, target Qubit) *Q {
	cnot := gate.ControlledZ(q.NumberOfBit(), Index(control...), target.Index())
	q.internal.Apply(cnot)
	return q
}

func (q *Q) CZ(control, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control}, target)
}

func (q *Q) CCZ(control0, control1, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control0, control1}, target)
}

func (q *Q) ControlledNot(control []Qubit, target Qubit) *Q {
	cnot := gate.ControlledNot(q.NumberOfBit(), Index(control...), target.Index())
	q.internal.Apply(cnot)
	return q
}

func (q *Q) CNOT(control, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control}, target)
}

func (q *Q) CCNOT(control0, control1, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control0, control1}, target)
}

func (q *Q) CCCNOT(control0, control1, control2, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control0, control1, control2}, target)
}

func (q *Q) Toffoli(control0, control1, target Qubit) *Q {
	return q.CCNOT(control0, control1, target)
}

func (q *Q) ConditionX(condition bool, qb ...Qubit) *Q {
	if condition {
		return q.X(qb...)
	}

	return q
}

func (q *Q) ConditionZ(condition bool, qb ...Qubit) *Q {
	if condition {
		return q.Z(qb...)
	}

	return q
}

func (q *Q) ControlledModExp2(a, j, N int, control Qubit, target []Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.CModExp2(n, a, j, N, control.Index(), Index(target...))
	q.internal.Apply(g)
	return q
}

func (q *Q) CModExp2(a, N int, control []Qubit, target []Qubit) *Q {
	for j := 0; j < len(control); j++ {
		q.ControlledModExp2(a, j, N, control[j], target)
	}

	return q
}

func (q *Q) Swap(qb ...Qubit) *Q {
	n := q.NumberOfBit()
	l := len(qb)

	for i := 0; i < l/2; i++ {
		q0, q1 := qb[i], qb[(l-1)-i]
		swap := gate.Swap(n, q0.Index(), q1.Index())
		q.internal.Apply(swap)
	}

	return q
}

func (q *Q) QFT(qb ...Qubit) *Q {
	l := len(qb)
	for i := 0; i < l; i++ {
		q.H(qb[i])

		k := 2
		for j := i + 1; j < l; j++ {
			q.CR(qb[j], qb[i], k)
			k++
		}
	}

	return q
}

func (q *Q) InverseQFT(qb ...Qubit) *Q {
	l := len(qb)
	for i := l - 1; i > -1; i-- {
		k := l - i
		for j := l - 1; j > i; j-- {
			q.CR(qb[j], qb[i], k)
			k--
		}

		q.H(qb[i])
	}

	return q
}

func (q *Q) InvQFT(qb ...Qubit) *Q {
	return q.InverseQFT(qb...)
}

func (q *Q) Clone() *Q {
	if q.internal == nil {
		return &Q{
			internal: nil,
			Seed:     q.Seed,
			Rand:     q.Rand,
		}
	}

	return &Q{
		internal: q.internal.Clone(),
		Seed:     q.internal.Seed,
		Rand:     q.internal.Rand,
	}
}

func (q *Q) String() string {
	return q.internal.String()
}

type State struct {
	Amplitude   complex128
	Probability float64
	Index       []int
	Binary      []string
}

func (s State) String() string {
	return fmt.Sprintf("%v%3v(% .4f% .4fi): %.4f", s.Binary, s.Index, real(s.Amplitude), imag(s.Amplitude), s.Probability)
}

func (q *Q) State(reg ...[]Qubit) []State {
	if len(reg) < 1 {
		qb := make([]Qubit, 0)
		for i := 0; i < q.NumberOfBit(); i++ {
			qb = append(qb, Qubit(i))
		}

		reg = append(reg, qb)
	}

	state := make([]State, 0)
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(q.NumberOfBit()), "s")
	for i, a := range q.internal.Amplitude() {
		if a == 0 {
			continue
		}
		if math.Abs(real(a)) < 1e-13 {
			a = complex(0, imag(a))
		}
		if math.Abs(imag(a)) < 1e-13 {
			a = complex(real(a), 0)
		}

		s := State{Amplitude: a, Probability: math.Pow(cmplx.Abs(a), 2)}
		bin := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))
		for _, r := range reg {
			var bb strings.Builder
			for _, b := range r {
				idx := b.Index()
				bb.WriteString(bin[idx : idx+1])
			}

			bbin := bb.String()
			bint, err := strconv.ParseInt(bbin, 2, 0)
			if err != nil {
				panic(fmt.Sprintf("parse int bin=%s, reg=%s", bin, bbin))
			}

			s.Index, s.Binary = append(s.Index, int(bint)), append(s.Binary, bbin)
		}

		state = append(state, s)
	}

	return state
}
