package matroid

type UniformMatroid struct {
	groundSet *Set
	n         int
}

func (u UniformMatroid) GroundSet() *Set {
	return u.groundSet
}

func (u UniformMatroid) Rank(s *Set) int {
	return min(s.Cardinality(), u.n)

}

func (u UniformMatroid) Independent(s *Set) bool {
	return s.Cardinality() == u.Rank(s)
}

func NewUniformMatroid(s *Set, n int) *UniformMatroid {
	return &UniformMatroid{
		groundSet: s,
		n:         n,
	}
}
