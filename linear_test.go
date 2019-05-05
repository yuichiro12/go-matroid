package matroid

import (
	"testing"
)

func TestRank(t *testing.T) {
	m := Matrix{
		NewUnweightedVector([]float64{1, 1, 1, 1}),
		NewUnweightedVector([]float64{0, 1, 1, 1}),
		NewUnweightedVector([]float64{0, 0, 1, 1}),
		NewUnweightedVector([]float64{0, 0, 0, 1}),
		NewUnweightedVector([]float64{0, 0, 0, 0}),
	}
	lm := NewLinearMatroid(m)

	if lm.Rank(lm.GroundSet()) != 4 {
		t.Errorf("rank mismatch. expected: 4, actual: %d", lm.Rank(lm.GroundSet()))
	}
	m = Matrix{
		NewUnweightedVector([]float64{1, 1, 1, 1, 1}),
		NewUnweightedVector([]float64{0, 0.5, 0.5, 0.5, 0.5}),
		NewUnweightedVector([]float64{0, 1, 1, 1, 1}),
		NewUnweightedVector([]float64{0, 2, 2, 2, 2}),
		NewUnweightedVector([]float64{0, 3, 3, 3, 3}),
		NewUnweightedVector([]float64{0, 4, 4, 4, 4}),
		NewUnweightedVector([]float64{0, 0, 0, 0, 0}),
	}
	lm = NewLinearMatroid(m)
	if lm.Rank(lm.GroundSet()) != 2 {
		t.Errorf("rank mismatch. expected: 2, actual: %d", lm.Rank(lm.GroundSet()))
	}
}
