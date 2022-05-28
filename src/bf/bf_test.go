package bf

import (
	"fmt"
	"math"
	"testing"
)

func TestFindNegativeCycle(t *testing.T) {
	edges := []Edge{
		{"a", "b"},
		{"b", "c"},
		{"c", "a"},
		{"c", "d"},
		{"d", "a"},
	}
	weight := map[string]float64{
		"a-b": -math.Log(1),
		"b-c": -math.Log(1),
		"c-a": -math.Log(1),
		"c-d": -math.Log(2),
		"d-a": -math.Log(1),
	}
	g := NewDirectedGraph(edges)
	f := func(f string, t string) float64 {
		return weight[f+"-"+t]
	}
	a, err := FindNegativeCycle(g, "a", f)
	fmt.Println(a, err)
}
