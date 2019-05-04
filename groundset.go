package matroid_intersection

import (
	"fmt"
	"strings"
)

// GroundSet implements Set
type GroundSet struct {
	set           map[string]Element
	groundSetType ElementType
}

type ElementType string

// Element is an atomic element of ground set.
// Be sure that Element is immutable.
type Element interface {
	GetType() ElementType
	Key() string
	Value() interface{}
}

func NewGroundSet(t ElementType, e ...Element) (*GroundSet, error) {
	gs := &GroundSet{
		set:           make(map[string]Element),
		groundSetType: t,
	}
	for _, el := range e {
		if err := gs.Add(el); err != nil {
			return nil, err
		}
	}
	return gs, nil
}

func EmptySet(t ElementType) *GroundSet {
	gs, _ := NewGroundSet(t)
	return gs
}

func errorTypeMismatch(t0, t1 ElementType) error {
	return fmt.Errorf("element type mismatch: %s and %s", t0, t1)
}

func errorGSTypeMismatch(t0, t1 ElementType) error {
	return fmt.Errorf("groundSet type mismatch: %s and %s", t0, t1)
}

func (gs *GroundSet) Add(e Element) error {
	if e.GetType() != gs.groundSetType {
		return errorTypeMismatch(gs.groundSetType, e.GetType())
	}
	gs.set[e.Key()] = e
	return nil
}

func (gs *GroundSet) Cardinality() int {
	var c int
	for range gs.set {
		c++
	}
	return c
}

func (gs *GroundSet) Clear() {
	gs.set = make(map[string]Element)
}

// Clone() is just a shallow clone. if your Element has field of slice, map or pointer,
// the reference to the original GroundSet still remains.
func (gs *GroundSet) Clone() *GroundSet {
	gs0 := EmptySet(gs.groundSetType)
	for e := range gs.Iter() {
		gs0.set[e.Key()] = e
	}
	return gs0
}

// TODO: implement deepClone

// Contains() returns true if given Elements are all in GroundSet.
// Note that this doesn't check groundSetType; only check Key() of given Elements.
func (gs *GroundSet) Contains(e ...Element) bool {
	for _, v := range e {
		if _, ok := gs.set[v.Key()]; !ok {
			return false
		}
	}
	return true
}

// Difference() subtracts `other` from `gs`.
// If your Element has field of slice, map or pointer, the reference to the original GroundSet still remains.
func (gs *GroundSet) Difference(other *GroundSet) (*GroundSet, error) {
	if gs.groundSetType != other.groundSetType {
		return nil, errorGSTypeMismatch(gs.groundSetType, other.groundSetType)
	}
	gs0 := EmptySet(gs.groundSetType)
	for e := range gs.Iter() {
		if !other.Contains(e) {
			_ = gs0.Add(e)
		}
	}
	return gs0, nil
}

func (gs *GroundSet) Equal(other *GroundSet) bool {
	return gs.IsSubset(other) && other.IsSubset(gs)
}

func (gs *GroundSet) Intersect(other *GroundSet) (*GroundSet, error) {
	if gs.groundSetType != other.groundSetType {
		return nil, errorGSTypeMismatch(gs.groundSetType, other.groundSetType)
	}
	gs0 := EmptySet(gs.groundSetType)
	for e := range gs.Iter() {
		if other.Contains(e) {
			_ = gs0.Add(e)
		}
	}
	return gs0, nil
}

func (gs *GroundSet) IsProperSubset(other *GroundSet) bool {
	return gs.IsSubset(other) && gs.Cardinality() < other.Cardinality()
}

func (gs *GroundSet) IsProperSuperset(other *GroundSet) bool {
	return other.IsProperSubset(gs)
}

func (gs *GroundSet) IsSubset(other *GroundSet) bool {
	for e := range gs.Iter() {
		if !other.Contains(e) {
			return false
		}
	}
	return true
}

func (gs *GroundSet) IsSuperset(other *GroundSet) bool {
	return other.IsSubset(gs)
}

func (gs *GroundSet) Each(f func(Element) bool) {
	for e := range gs.Iter() {
		if !f(e) {
			break
		}
	}
}

func (gs *GroundSet) Iter() <-chan Element {
	ch := make(chan Element)
	go func() {
		for _, e := range gs.set {
			ch <- e
		}
		close(ch)
	}()

	return ch
}

// Remove() removes given element from GroundSet.
// Note that doesn't check groundSetType; only check Key() of given Elements.
func (gs *GroundSet) Remove(e Element) {
	if _, ok := gs.set[e.Key()]; ok {
		delete(gs.set, e.Key())
	}
}

func (gs *GroundSet) SymmetricDifference(other *GroundSet) (*GroundSet, error) {
	if gs.groundSetType != other.groundSetType {
		return nil, errorGSTypeMismatch(gs.groundSetType, other.groundSetType)
	}
	gs1, _ := gs.Difference(other)
	gs2, _ := other.Difference(gs)
	gs0, _ := gs1.Union(gs2)
	return gs0, nil
}

func (gs *GroundSet) Union(other *GroundSet) (*GroundSet, error) {
	if gs.groundSetType != other.groundSetType {
		return nil, errorGSTypeMismatch(gs.groundSetType, other.groundSetType)
	}
	gs0 := EmptySet(gs.groundSetType)
	diff, _ := gs.Difference(other)
	for e := range diff.Iter() {
		_ = gs0.Add(e)
	}
	for e := range other.Iter() {
		_ = gs0.Add(e)
	}
	return gs0, nil
}

func (gs *GroundSet) String() string {
	var s []string
	for e := range gs.Iter() {
		s = append(s, " "+e.Key())
	}
	return "GroundSet{\n" + strings.Join(s, "\n") + "\n}"
}

// Pop() removes and returns an arbitrary element from GroundSet.
// if the GroundSet is empty, Pop() returns nil.
func (gs *GroundSet) Pop() Element {
	for e := range gs.Iter() {
		return e
	}
	return nil
}

// Choose() returns an arbitrary element that makes given callback to be true.
// if none of elements in GroundSet satisfy given callback, Choose() returns nil.
func (gs *GroundSet) Choose(f func(Element) bool) Element {
	for e := range gs.Iter() {
		if f(e) {
			return e
		}
	}
	return nil
}

// CondSubset() returns a new subset consisting of elements that make given callback to be true
func (gs *GroundSet) CondSubset(f func(Element) bool) *GroundSet {
	gs0 := EmptySet(gs.groundSetType)
	for e := range gs.Iter() {
		if f(e) {
			_ = gs0.Add(e)
		}
	}
	return gs0
}

func (gs *GroundSet) IsEmpty() bool {
	return gs.Cardinality() == 0
}

func (gs *GroundSet) ToSlice() []Element {
	var elms []Element
	for e := range gs.Iter() {
		elms = append(elms, e)
	}
	return elms
}
