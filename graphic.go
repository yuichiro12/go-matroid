package matroid

import "fmt"

const ArcType ElementType = "ARC"
const VertexType ElementType = "VERTEX"

type Arc struct {
	Start *Vertex
	End   *Vertex
	W     float64
	Id    string
}

func (e *Arc) From() *Vertex {
	return e.Start
}

func (e *Arc) To() *Vertex {
	return e.End
}

func (e *Arc) GetType() ElementType {
	return ArcType
}

func (e *Arc) Key() string {
	return fmt.Sprintf("%s:d(%s,%s)", e.Id, e.From().Key(), e.To().Key())
}

func (e *Arc) Value() interface{} {
	return [2]*Vertex{e.Start, e.End}
}

func (e *Arc) Weight() float64 {
	return e.W
}

type Vertex struct {
	ID string
	W  float64
}

func (v *Vertex) GetType() ElementType {
	return VertexType
}

func (v *Vertex) Key() string {
	return fmt.Sprintf("%d", v.ID)
}

func (v *Vertex) Value() interface{} {
	return v.ID
}

func (v *Vertex) Weight() float64 {
	return v.W
}

type WeightedDigraph struct {
	// V is set of Vertices
	V *Set
	// A is set of Arcs
	A *Set
}

type Weighted

// func ToAdjacentMatrix()
// func ToIncidenceMatrix()
