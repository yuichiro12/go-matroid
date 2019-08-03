package matroid

import "errors"

// `partition` field has other two fields info, but we preserve them for convenient
type PartitionMatroid struct {
	groundSet  *Set
	partitions []Partition
}

// partition of ground set with intersect threshold n
// n must be non-negative
type Partition struct {
	set *Set
	n   int
}

func UnionAllPartitions(p []Partition) (*Set, error) {
	s := EmptySet(p[0].set.GetType())
	var c int
	var r int
	for _, pp := range p {
		s.Union(pp.set)
		c += pp.set.Cardinality()
		r += pp.n
	}
	if s.Cardinality() != c {
		return nil, errors.New("given sets are not disjoint")
	}
	return s, nil
}

func (p PartitionMatroid) GroundSet() *Set {
	return p.groundSet
}

func (p PartitionMatroid) Rank(s *Set) int {
	var r int
	for _, pp := range p.partitions {
		r += min(s.Intersect(pp.set).Cardinality(), pp.n)
	}
	return r
}

func (p PartitionMatroid) Independent(s *Set) bool {
	return s.Cardinality() == p.Rank(s)
}

func NewPartitionMatroid(s ...*Set) (*PartitionMatroid, error) {
	var p []Partition
	for _, ss := range s {
		p = append(p, Partition{
			set: ss,
			n:   1,
		})
	}
	return NewGeneralizedPartitionMatroid(p)
}

func NewGeneralizedPartitionMatroid(p []Partition) (*PartitionMatroid, error) {
	s, err := UnionAllPartitions(p)
	return &PartitionMatroid{
		groundSet:  s,
		partitions: p,
	}, err
}
