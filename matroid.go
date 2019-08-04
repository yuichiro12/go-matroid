package matroid

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/graph"

	"gonum.org/v1/gonum/graph/simple"
)

type Matroid interface {
	// GroundSet() returns GroundSet of matroid
	GroundSet() *Set
	// Rank() is rank oracle of the matroid.
	// Make sure that input Set must be a subset of GroundSet.
	Rank(*Set) int
	// Independent() returns true if given Set is independent set of matroid.
	// Make sure that input Set must be a subset of GroundSet.
	// This is easily implemented with Rank() function. For example, see Matroid implementors
	// in this package.
	Independent(*Set) bool
}

type sorter []Element

func (s sorter) Len() int {
	return len(s)
}

func (s sorter) Less(i, j int) bool {
	return s[i].Weight() < s[j].Weight()
}

func (s sorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Intersection() returns maximal matroid intersection of input two matroids.
func Intersection(m1, m2 Matroid) (*Set, error) {
	if !(m1.GroundSet().GetType() == m2.GroundSet().GetType()) {
		return nil, fmt.Errorf("incomparable setTypes: %s and %s",
			m1.GroundSet().GetType(), m2.GroundSet().GetType())
	}
	if !m1.GroundSet().Equal(m2.GroundSet()) {
		return nil, fmt.Errorf("inequal GroundSets")
	}
	gs := m1.GroundSet()
	s := EmptySet(gs.GetType())

	for e := range gs.Iter() {
		s.Add(e)
		if !(m1.Independent(s) && m2.Independent(s)) {
			s.Remove(e)
		}
	}

	for {
		c, _ := gs.Complement(s)
		d := generateMatroidIntersectionBipartiteDigraph(s, c, m1, m2)
		p := findShortestPath(d)
		if p == nil {
			break
		}
		swapAlongPath(s, p)
	}
	return s, nil
}

func generateMatroidIntersectionBipartiteDigraph(s, c *Set, m1, m2 Matroid) *simple.WeightedDirectedGraph {
	nodes := getKeyToNodeMap(s, c)
	d := simple.NewWeightedDirectedGraph(0, math.Inf(1))
	for _, v := range nodes {
		d.AddNode(v)
	}
	s0 := s.Clone()

	for f := range c.Iter() {
		s0.Add(f)
		if m1.Independent(s0) {
			nodes[f.Key()].isSink = true
		}
		if m2.Independent(s0) {
			nodes[f.Key()].isSource = true
		}
		s0.Remove(f)
	}
	for e := range s.Iter() {
		for f := range c.Iter() {
			s0.Swap(f, e)
			if m1.Independent(s0) {
				d.SetWeightedEdge(&weightedEdge{tail: nodes[e.Key()], head: nodes[f.Key()]})
			}
			if m2.Independent(s0) {
				d.SetWeightedEdge(&weightedEdge{tail: nodes[f.Key()], head: nodes[e.Key()]})
			}
			s0.Swap(e, f)
		}
	}
	return d
}

// augment s
func swapAlongPath(s *Set, p []graph.Node) {
	for i, v := range p {
		if i%2 == 0 {
			s.Add(v.(*node).element)
		} else {
			s.Remove(v.(*node).element)
		}
	}
}

func getKeyToNodeMap(s, c *Set) map[string]*node {
	m := make(map[string]*node)
	var idx int64
	for e := range s.Iter() {
		m[e.Key()] = &node{id: idx, element: e}
		idx++
	}
	for e := range c.Iter() {
		m[e.Key()] = &node{id: idx, element: e}
		idx++
	}
	return m
}

// GetBaseOf() returns an arbitrary base of input matroid.
func GetBaseOf(m Matroid) *Set {
	set := EmptySet(m.GroundSet().GetType())

	s := m.GroundSet().ToSlice()
	for i := 0; i < len(s); i++ {
		set.Add(s[i])
		if !m.Independent(set) {
			set.Remove(s[i])
		}
	}
	return set
}

// GetMaximalBaseOf() returns maximal base of input matroid.
func GetMaximalBaseOf(m Matroid) *Set {
	set := EmptySet(m.GroundSet().GetType())
	var s sorter
	s = m.GroundSet().ToSlice()
	sort.Sort(s)
	for i := 0; i < len(s); i++ {
		set.Add(s[i])
		if !m.Independent(set) {
			set.Remove(s[i])
		}
	}
	return set
}

// Dual() returns dual matroid of input matroid.
func Dual(m Matroid) Matroid {
	return &dualMatroid{
		groundSet: m.GroundSet(),
		r:         m.Rank,
	}
}

type dualMatroid struct {
	groundSet *Set
	// rank function of original matroid
	r func(*Set) int
}

func (dm *dualMatroid) GroundSet() *Set {
	return dm.groundSet
}

func (dm *dualMatroid) Rank(s *Set) int {
	c, err := dm.GroundSet().Complement(s)
	// make sure that input is the subset of the GroundSet
	if err != nil {
		panic(err)
	}
	return dm.r(c) + s.Cardinality() - dm.r(dm.GroundSet())
}

func (dm *dualMatroid) Independent(s *Set) bool {
	return s.Cardinality() == dm.Rank(s)
}
