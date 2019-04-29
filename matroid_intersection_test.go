package matroid_intersection

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestLinearMatroidRank(t *testing.T) {
	m := mat.NewDense(4, 4, nil)
	m.SetRow(0, []float64{1, 1, 1, 1})
	m.SetRow(1, []float64{0, 1, 1, 1})
	m.SetRow(2, []float64{0, 0, 1, 1})
	m.SetRow(3, []float64{0, 0, 0, 1})
	lm := &LinearMatroid{
		E: m,
	}
	if lm.rank() != 4 {
		t.Fatalf("expected: 4, actual: %d", lm.rank())
	}
	m.SetRow(2, []float64{0, 1, 1, 1})
	m.SetRow(3, []float64{0, 1, 1, 1})
	if lm.rank() != 2 {
		t.Fatalf("expected: 2, actual: %d", lm.rank())
	}
}
