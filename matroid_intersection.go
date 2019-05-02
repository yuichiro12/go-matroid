package matroid_intersection

import (
	mapset "github.com/deckarep/golang-set"
)

type Matroid interface {
	// rank oracle of the matroid
	Rank() int
	// validate if element e can be added to the GroundSet
	Validate(e interface{}) bool
}

type PartitionMatroid struct {
	E mapset.Set
}

func Intersection(m1, m2 Matroid) {
	mapset.NewSet()
}

func isEqualGroundSet(m1, m2 Matroid) {
}
