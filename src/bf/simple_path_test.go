package bf

import (
	"log"
	"testing"
)

func TestAllSimplePaths(t *testing.T) {
	edges := []Edge{
		{"a", "b1"},
		{"a", "c1"},
		{"b1", "b2"},
		{"b3", "z"},
		{"c1", "c2"},
		{"c2", "z"},
		{"c1", "z"},
		{"c1", "c1-1"},
		{"c1-1", "z"},
		// {"c1-1", "c1-1"},
		{"c1-1", "c1"},
	}
	g := NewDirectedGraph(edges)
	ch := AllSimplePaths(g, "a", "z", 3)
	for rs := range ch {
		log.Printf("%#v", rs)
	}
}
