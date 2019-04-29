package matroid_intersection

import "gonum.org/v1/gonum/mat"

type matroid interface {
	rank() int
}

type LinearMatroid struct {
	// E is ground set of this matroid
	E mat.Matrix
}

func rank(m mat.Matrix) int {
	svd := new(mat.SVD)

	svd.Factorize(m, mat.SVDNone)
	svs := svd.Values(nil)
	var count int
	for _, v := range svs {
		if v != 0 {
			count++
		}
	}
	return count
}

func (lm *LinearMatroid) rank() int {
	return rank(lm.E)
}
