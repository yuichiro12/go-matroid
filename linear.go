package matroid_intersection

import (
	"fmt"
	"strings"

	"gonum.org/v1/gonum/mat"
)

const epsilon = 1e-10
const VectorType ElementType = "VECTOR"

type LinearMatroid struct {
	*GroundSet
}

// Vector implements Element
type Vector struct {
	Value  []float64
	Weight float64
}

func (v Vector) Key() string {
	var s []string
	for i := 0; i < len(v.Value); i++ {
		s = append(s, fmt.Sprintf("%f", v.Value[i]))
	}
	return "(" + strings.Join(s, ",") + ")"
}

func (v Vector) GetType() ElementType {
	return VectorType
}

func NewUnweightedVector(v []float64) Vector {
	return NewWeightedVector(0, v)
}

func NewWeightedVector(w float64, v []float64) Vector {
	return Vector{
		Value:  v,
		Weight: w,
	}
}

// Matrix implements mat.Matrix
type Matrix []Vector

// At() returns the value of a matrix element at row i, column j
func (m Matrix) At(i, j int) float64 {
	return m[i].Value[j]
}

//
func (m Matrix) Dims() (r, c int) {
	if len(m) == 0 {
		return 0, 0
	} else {
		return len(m), len(m[0].Value)
	}
}

func (m Matrix) T() mat.Matrix {
	return mat.Transpose{m}
}

// each ROW of input matrix will be an element of groundSet (not COLUMNS)
func NewLinearMatroid(m Matrix) *LinearMatroid {
	gs := EmptySet(VectorType)
	for _, e := range m {
		_ = gs.Add(e)
	}
	return &LinearMatroid{
		GroundSet: gs,
	}
}

// the order of rows is not idempotent because the set has no order
func (lm *LinearMatroid) AsMatrix() Matrix {
	var m Matrix
	for e := range lm.Iter() {
		m = append(m, e.(Vector))
	}
	return m
}

func rank(m mat.Matrix) int {
	svd := new(mat.SVD)
	svd.Factorize(m, mat.SVDNone)
	svs := svd.Values(nil)
	var count int
	for _, v := range svs {
		if v > epsilon {
			count++
		}
	}
	return count
}

func (lm *LinearMatroid) Rank() int {
	return rank(lm.AsMatrix())
}
