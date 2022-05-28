package bf

import (
	"math"
	"reflect"
	"testing"
)

func TestFindNegativeCycle1(t *testing.T) {
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
	if err != nil {
		t.Fatal("err", err)
	}
	if !reflect.DeepEqual(StartFrom(a, "a"), []string{"a", "b", "c", "d", "a"}) {
		t.Fatal("route error", a)
	}
}

func TestFindNegativeCycle2(t *testing.T) {
	edges := []Edge{
		{"a", "b"},
		{"b", "c"},
		{"c", "a"},
	}
	weight := map[string]float64{
		"a-b": -math.Log(1),
		"b-c": -math.Log(1),
		"c-a": -math.Log(1),
	}
	g := NewDirectedGraph(edges)
	f := func(f string, t string) float64 {
		return weight[f+"-"+t]
	}
	_, err := FindNegativeCycle(g, "a", f)
	if err == nil {
		t.Fatal("there is no negative cycle")
	}
}

func TestFindNegativeCycle3(t *testing.T) {
	edges := []Edge{
		{"a", "b"},
		{"b", "c"},
	}
	weight := map[string]float64{
		"a-b": -math.Log(1),
		"b-c": -math.Log(1),
	}
	g := NewDirectedGraph(edges)
	f := func(f string, t string) float64 {
		return weight[f+"-"+t]
	}
	_, err := FindNegativeCycle(g, "a", f)
	if err == nil {
		t.Fatal("there is no cycle")
	}
}

func TestFindNegativeCycle4(t *testing.T) {
	edges := []Edge{
		{"a", "a"},
	}
	weight := map[string]float64{
		"a-a": -1,
	}
	g := NewDirectedGraph(edges)
	f := func(f string, t string) float64 {
		return weight[f+"-"+t]
	}
	path, err := FindNegativeCycle(g, "a", f)
	if err != nil {
		t.Fatal("there is cycle")
	}
	if !reflect.DeepEqual(StartFrom(path, "a"), []string{"a", "a"}) {
		t.Fatal("route error", path)
	}
}

func TestFindNegativeCycle5(t *testing.T) {
	edges := []Edge{
		{"a", "a"},
	}
	weight := map[string]float64{
		"a-a": 1,
	}
	g := NewDirectedGraph(edges)
	f := func(f string, t string) float64 {
		return weight[f+"-"+t]
	}
	_, err := FindNegativeCycle(g, "a", f)
	if err == nil {
		t.Fatal("there is no cycle")
	}
}

func TestFindNegativeCycle6(t *testing.T) {
	edges := []Edge{
		{"a", "b"},
	}
	weight := map[string]float64{
		"a-b": -math.Log(1),
	}
	g := NewDirectedGraph(edges)
	f := func(f string, t string) float64 {
		return weight[f+"-"+t]
	}
	_, err := FindNegativeCycle(g, "a", f)
	if err == nil {
		t.Fatal("there is no cycle")
	}
}

func TestFindNegativeCycle7(t *testing.T) {
	edges := []Edge{
		{"a", "b"},
		{"b", "c"},
		{"b", "d"},
		{"d", "a"},
	}
	weight := map[string]float64{
		"a-b": -math.Log(1),
		"b-c": -math.Log(1),
		"b-d": -math.Log(1),
		"d-a": -math.Log(1),
	}
	g := NewDirectedGraph(edges)
	f := func(f string, t string) float64 {
		return weight[f+"-"+t]
	}
	_, err := FindNegativeCycle(g, "a", f)
	if err == nil {
		t.Fatal("there is no cycle")
	}
}
