package matroid_intersection

import "testing"

func TestRank(t *testing.T) {
	m := Matrix{
		NewUnweightedVector([]float64{1, 1, 1, 1}),
		NewUnweightedVector([]float64{0, 1, 1, 1}),
		NewUnweightedVector([]float64{0, 0, 1, 1}),
		NewUnweightedVector([]float64{0, 0, 0, 1}),
		NewUnweightedVector([]float64{0, 0, 0, 0}),
	}
	lm := NewLinearMatroid(m)
	if lm.Rank() != 4 {
		t.Errorf("rank mismatch. expected: 4, actual: %d", lm.Rank())
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
	if lm.Rank() != 2 {
		t.Errorf("rank mismatch. expected: 2, actual: %d", lm.Rank())
	}
}
