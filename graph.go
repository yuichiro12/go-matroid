package matroid

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
)

const ArcType ElementType = "ARC"
const VertexType ElementType = "VERTEX"

type Edge struct {
	Start *Node
	End   *Node
	W     float64
	Label string
}

func (e *Edge) From() graph.Node {
	return e.Start
}

func (e *Edge) To() graph.Node {
	panic("implement me")
}

func (e *Edge) ReversedEdge() graph.Edge {
	panic("implement me")
}

func (e *Edge) GetType() ElementType {
	return ArcType
}

func (e *Edge) Key() string {
	return fmt.Sprintf("%s:d(%d,%d)", e.Label, e.From().ID(), e.To().ID())
}

func (e *Edge) Value() interface{} {
	return [2]*Node{e.Start, e.End}
}

func (e *Edge) Weight() float64 {
	return e.W
}

type Node struct {
	Id int64
	W  float64
}

func (v *Node) ID() int64 {
	return v.Id
}

func (v *Node) GetType() ElementType {
	return VertexType
}

func (v *Node) Key() string {
	return fmt.Sprintf("%d", v.Id)
}

func (v *Node) Value() interface{} {
	return v.Id
}

func (v *Node) Weight() float64 {
	return v.W
}

// WeightedDigraph implements graph.WeightedDirected
type WeightedDigraph struct {
	// V is set of Vertices
	V *Set
	// A is set of Arcs
	A *Set
}

func (d WeightedDigraph) Node(id int64) graph.Node {
	panic("implement me")
}

func (d WeightedDigraph) Nodes() graph.Nodes {
	panic("implement me")
}

func (d WeightedDigraph) From(id int64) graph.Nodes {
	panic("implement me")
}

func (d WeightedDigraph) HasEdgeBetween(xid, yid int64) bool {
	panic("implement me")
}

func (d WeightedDigraph) Edge(uid, vid int64) graph.Edge {
	panic("implement me")
}

func (d WeightedDigraph) WeightedEdge(uid, vid int64) graph.WeightedEdge {
	panic("implement me")
}

func (d WeightedDigraph) Weight(xid, yid int64) (w float64, ok bool) {
	panic("implement me")
}

func (d WeightedDigraph) HasEdgeFromTo(uid, vid int64) bool {
	panic("implement me")
}

func (d WeightedDigraph) To(id int64) graph.Nodes {
	panic("implement me")
}
