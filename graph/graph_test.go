package graph_test

import (
	"testing"

	"go.jlucktay.dev/golang-workbench/graph"
)

var g graph.ItemGraph

func fillGraph() {
	nA := graph.NewNode("A")
	nB := graph.NewNode("B")
	nC := graph.NewNode("C")
	nD := graph.NewNode("D")
	nE := graph.NewNode("E")
	nF := graph.NewNode("F")
	g.AddNode(&nA)
	g.AddNode(&nB)
	g.AddNode(&nC)
	g.AddNode(&nD)
	g.AddNode(&nE)
	g.AddNode(&nF)

	g.AddEdge(&nA, &nB)
	g.AddEdge(&nA, &nC)
	g.AddEdge(&nB, &nE)
	g.AddEdge(&nC, &nE)
	g.AddEdge(&nE, &nF)
	g.AddEdge(&nD, &nA)
}

func TestAdd(t *testing.T) {
	fillGraph()
	g.String()
}
