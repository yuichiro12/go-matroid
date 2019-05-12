package matroid

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/multi"
)

const ArcType ElementType = "ARC"
const VertexType ElementType = "VERTEX"

type Arc struct {
	Tail *Vertex
	Head *Vertex
	W    float64
	Id   int64
}

func AsArc(l graph.Line) *Arc {
	return l.(*Arc)
}

func (l *Arc) GetType() ElementType {
	return ArcType
}

func (l *Arc) Key() string {
	return fmt.Sprintf("%d:(%d,%d)")
}

func (l *Arc) Value() interface{} {
	return l.Id
}

func (l *Arc) Weight() float64 {
	return l.W
}

func (l *Arc) From() graph.Node {
	return l.Tail
}

func (l *Arc) To() graph.Node {
	return l.Head
}

func (l *Arc) ReversedLine() graph.Line {
	return &Arc{
		Tail: l.Head,
		Head: l.Tail,
		W:    l.W,
		Id:   l.Id,
	}
}

func (l *Arc) ReversedNewLine(id int64) graph.Line {
	return &Arc{
		Tail: l.Head,
		Head: l.Tail,
		W:    l.W,
		Id:   id,
	}
}

func (l *Arc) ID() int64 {
	return l.Id
}

type Vertex struct {
	Id int64
	W  float64
}

func AsVertex(n graph.Node) *Vertex {
	return n.(*Vertex)
}

func (v *Vertex) GetType() ElementType {
	return VertexType
}

func (v *Vertex) Key() string {
	return fmt.Sprintf("%d", v.Id)
}

func (v *Vertex) Value() interface{} {
	return v.ID()
}

func (v *Vertex) Weight() float64 {
	return v.W
}

func (v Vertex) ID() int64 {
	return v.Id
}

type WeightedDigraph struct {
	*multi.WeightedDirectedGraph
	A *Set
	V *Set
}

func (d *WeightedDigraph) AddVertex(v *Vertex) {
	d.V.Add(v)
	d.AddNode(v)
}

func (d *WeightedDigraph) RemoveVertex(v *Vertex) {
	d.V.Remove(v)
	d.RemoveNode(v.Id)
}

func (d *WeightedDigraph) AddArc(a *Arc) {
	d.A.Add(a)
	d.SetWeightedLine(a)
}

func (d *WeightedDigraph) RemoveArc(a *Arc) {
	d.A.Remove(a)
	d.RemoveLine(a.Tail.Id, a.Head.Id, a.Id)
}

func NewWeightedDigraph() *WeightedDigraph {
	return &WeightedDigraph{
		WeightedDirectedGraph: multi.NewWeightedDirectedGraph(),
		A:                     EmptySet(ArcType),
		V:                     EmptySet(VertexType),
	}
}
