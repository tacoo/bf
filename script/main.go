package main

import (
	"fmt"

	"github.com/tacoo/bf/src/bf"
)

func main() {
	edges := []bf.Edge{
		{"a", "b"},
		{"b", "c"},
		{"c", "a"},
	}
	fmt.Printf("%#v\n", bf.NewDirectedGraph(edges))
}
