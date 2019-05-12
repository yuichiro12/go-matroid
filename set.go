package matroid

import (
	"errors"
	"fmt"
	"strings"
)

type Set struct {
	set     map[string]Element
	setType ElementType
}

// ElementType is type of element.
// This is mainly used for two validation:
// 1. checking if an element can be added to the set
// 2. checking if two Sets are comparable
type ElementType string

// Element is an atomic element of the Set.
// Be sure that Element is immutable.
type Element interface {
	// GetType() returns the type of this element
	GetType() ElementType
	// Key() returns unique and idempotent string related with Value().
	// This is used for the key of the map since Set is implemented with map.
	// Thanks for this Key(), you can add any kind of data (map, slice, func...)
	// no matter how they are unhashable.
	Key() string
	// Value() returns the value of this element
	Value() interface{}
	// Weight() returns the weight of this element
	Weight() float64
}

func NewSet(t ElementType, e ...Element) *Set {
	s := &Set{
		set:     make(map[string]Element),
		setType: t,
	}
	for _, el := range e {
		s.Add(el)
	}
	return s
}

func (s *Set) GetType() ElementType {
	return s.setType
}

func EmptySet(t ElementType) *Set {
	return NewSet(t)
}

func typeMismatchPanic(t0, t1 ElementType) {
	panic(fmt.Sprintf("ElementType mismatch: %s and %s", t0, t1))
}

// Add() returns true if element is added; otherwise false
func (s *Set) Add(e Element) bool {
	if e.GetType() != s.setType {
		typeMismatchPanic(s.setType, e.GetType())
	}
	if _, ok := s.set[e.Key()]; !ok {
		s.set[e.Key()] = e
		return true
	}
	return false
}

func (s *Set) Cardinality() int {
	var c int
	for range s.set {
		c++
	}
	return c
}

func (s *Set) Clear() {
	s.set = make(map[string]Element)
}

// Clone() is just a shallow clone. if your Element has field of slice, map or pointer,
// the reference to the original Set still remains.
// Basically this is no problem because Elements are immutable.
func (s *Set) Clone() *Set {
	s0 := EmptySet(s.setType)
	for e := range s.Iter() {
		s0.set[e.Key()] = e
	}
	return s0
}

// Contains() returns true if given Elements are all in Set.
// Note that this doesn't check setType; only check Key() of given Elements.
func (s *Set) Contains(e ...Element) bool {
	for _, v := range e {
		if _, ok := s.set[v.Key()]; !ok {
			return false
		}
	}
	return true
}

// Difference() subtracts `other` from `s`.
// If your Element has field of slice, map or pointer, the reference to the original Set still remains.
func (s *Set) Difference(other *Set) *Set {
	if s.setType != other.setType {
		typeMismatchPanic(s.setType, other.setType)
	}
	s0 := EmptySet(s.setType)
	for e := range s.Iter() {
		if !other.Contains(e) {
			s0.Add(e)
		}
	}
	return s0
}

func (s *Set) Equal(other *Set) bool {
	return s.IsSubsetOf(other) && other.IsSubsetOf(s)
}

func (s *Set) Intersect(other *Set) *Set {
	if s.setType != other.setType {
		typeMismatchPanic(s.setType, other.setType)
	}
	s0 := EmptySet(s.setType)
	for e := range s.Iter() {
		if other.Contains(e) {
			s0.Add(e)
		}
	}
	return s0
}

func (s *Set) IsProperSubsetOf(other *Set) bool {
	return s.IsSubsetOf(other) && s.Cardinality() < other.Cardinality()
}

func (s *Set) IsProperSupersetOf(other *Set) bool {
	return other.IsProperSubsetOf(s)
}

func (s *Set) IsSubsetOf(other *Set) bool {
	for e := range s.Iter() {
		if !other.Contains(e) {
			return false
		}
	}
	return true
}

func (s *Set) IsSupersetOf(other *Set) bool {
	return other.IsSubsetOf(s)
}

func (s *Set) Each(f func(Element) bool) {
	for e := range s.Iter() {
		if !f(e) {
			break
		}
	}
}

func (s *Set) Iter() <-chan Element {
	ch := make(chan Element)
	go func() {
		for _, e := range s.set {
			ch <- e
		}
		close(ch)
	}()

	return ch
}

// Remove() removes given element from Set.
// Note that doesn't check setType; only check Key() of given Elements.
func (s *Set) Remove(e Element) {
	delete(s.set, e.Key())
}

func (s *Set) SymmetricDifference(other *Set) *Set {
	if s.setType != other.setType {
		typeMismatchPanic(s.setType, other.setType)
	}
	return s.Difference(other).Union(other.Difference(s))
}

func (s *Set) Union(other *Set) *Set {
	if s.setType != other.setType {
		typeMismatchPanic(s.setType, other.setType)
	}
	s0 := EmptySet(s.setType)
	diff := s.Difference(other)
	for e := range diff.Iter() {
		s0.Add(e)
	}
	for e := range other.Iter() {
		s0.Add(e)
	}
	return s0
}

func (s *Set) String() string {
	var sl []string
	for e := range s.Iter() {
		sl = append(sl, " "+e.Key())
	}
	return "Set{\n" + strings.Join(sl, "\n") + "\n}"
}

// Pop() removes and returns an arbitrary element from Set.
// if the Set is empty, Pop() returns nil.
func (s *Set) Pop() Element {
	for e := range s.Iter() {
		return e
	}
	return nil
}

// Choose() returns an arbitrary element that makes given callback to be true.
// if none of elements in Set satisfy given callback, Choose() returns nil.
func (s *Set) Choose(f func(Element) bool) Element {
	for e := range s.Iter() {
		if f(e) {
			return e
		}
	}
	return nil
}

// CondSubset() returns a new subset consisting of elements that make given callback to be true
func (s *Set) CondSubset(f func(Element) bool) *Set {
	s0 := EmptySet(s.setType)
	for e := range s.Iter() {
		if f(e) {
			s0.Add(e)
		}
	}
	return s0
}

func (s *Set) IsEmpty() bool {
	return s.Cardinality() == 0
}

func (s *Set) ToSlice() []Element {
	var elms []Element
	for e := range s.Iter() {
		elms = append(elms, e)
	}
	return elms
}

func (s *Set) Complement(subset *Set) (*Set, error) {
	if !s.IsSupersetOf(subset) {
		return nil, errors.New("Complement(): input Set is not subset of the receiver Set")
	}
	return s.Difference(subset), nil
}
