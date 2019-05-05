package matroid

import (
	"fmt"
	"strings"

	"gonum.org/v1/gonum/mat"
)

var Epsilon = 1e-10

const VectorType ElementType = "VECTOR"

type LinearMatroid struct {
	groundSet *Set
}

func (lm *LinearMatroid) GroundSet() *Set {
	return lm.groundSet
}

func (lm *LinearMatroid) Rank(s *Set) int {
	return rank(lm.GetMatrixOf(s))
}

func (lm *LinearMatroid) Independent(s *Set) bool {
	return s.Cardinality() == lm.Rank(s)
}

// Vector implements Element
type Vector struct {
	V []float64
	W float64
}

func (v Vector) Key() string {
	var s []string
	for i := 0; i < len(v.V); i++ {
		s = append(s, fmt.Sprintf("%f", v.V[i]))
	}
	return "(" + strings.Join(s, ",") + ")"
}

func (v Vector) GetType() ElementType {
	return VectorType
}

func (v Vector) Value() interface{} {
	return v.V
}

func (w Vector) Weight() float64 {
	return w.W
}

func NewUnweightedVector(v []float64) Vector {
	return NewWeightedVector(0, v)
}

func NewWeightedVector(w float64, v []float64) Vector {
	return Vector{
		V: v,
		W: w,
	}
}

// Matrix implements mat.Matrix
type Matrix []Vector

func (m Matrix) At(i, j int) float64 {
	return m[i].V[j]
}

func (m Matrix) Dims() (r, c int) {
	if len(m) == 0 {
		return 0, 0
	} else {
		return len(m), len(m[0].V)
	}
}

func (m Matrix) T() mat.Matrix {
	return mat.Transpose{m}
}

// Each Vector of the input Matrix will be an element of the GroundSet.
// Be sure that each Vector is unique, that is, has a unique Key(), because GroundSet is a Set;
// otherwise duplicate Vectors are omitted except first one.
func NewLinearMatroid(m Matrix) *LinearMatroid {
	gs := EmptySet(VectorType)
	for _, e := range m {
		_ = gs.Add(e)
	}
	return &LinearMatroid{
		groundSet: gs,
	}
}

// GetMatrixOf() returns Matrix form of input Set.
// The input must be the subset of the GroundSet.
// the order of rows is not idempotent because the set has no order
func (lm *LinearMatroid) GetMatrixOf(s *Set) Matrix {
	var m Matrix
	for e := range s.Iter() {
		m = append(m, e.(Vector))
	}
	return m
}

// rank() returns rank of input matrix using singular value decomposition.
func rank(m mat.Matrix) int {
	svd := new(mat.SVD)
	svd.Factorize(m, mat.SVDNone)
	svs := svd.Values(nil)
	var count int
	for _, v := range svs {
		if v > Epsilon {
			count++
		}
	}
	return count
}
