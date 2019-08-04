package matroid

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

type node struct {
	id       int64
	element  Element
	isSource bool
	isSink   bool
}

func (n *node) ID() int64 {
	return n.id
}

type weightedEdge struct {
	tail *node
	head *node
}

func (w *weightedEdge) From() graph.Node {
	return w.tail
}

func (w *weightedEdge) To() graph.Node {
	return w.head
}

func (w *weightedEdge) ReversedEdge() graph.Edge {
	return &weightedEdge{
		tail: w.head,
		head: w.tail,
	}
}

func (w *weightedEdge) Weight() float64 {
	return w.head.element.Weight() - w.tail.element.Weight()
}

// findShortestPath() applies BFS to the input graph and return first found path from source to sink
func findShortestPath(d *simple.WeightedDirectedGraph) []graph.Node {
	paths := make(map[int64][]graph.Node)
	var bfs traverse.BreadthFirst
	var from graph.Node
	bfs.Traverse = func(edge graph.Edge) bool {
		from = edge.From()
		return !bfs.Visited(edge.To())
	}
	bfs.Visit = func(i graph.Node) {
		paths[i.ID()] = append(paths[from.ID()], i)
	}
	it := d.Nodes()
	for it.Next() {
		if !it.Node().(*node).isSource {
			continue
		}
		paths[it.Node().ID()] = []graph.Node{it.Node()}
		end := bfs.Walk(d, it.Node(), func(n graph.Node, depth int) bool {
			return n.(*node).isSink
		})
		if end != nil {
			return paths[end.ID()]
		}
		paths = make(map[int64][]graph.Node)
	}
	return nil
}
