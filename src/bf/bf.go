package bf

import (
	"errors"
	"math"
)

type Weight func(from string, to string) (weight float64)

type Edge [2]string

type DirectedGraph struct {
	adjacency map[string]map[string]struct{}
}

func NewDirectedGraph(edges []Edge) *DirectedGraph {
	g := DirectedGraph{
		adjacency: make(map[string]map[string]struct{}),
	}
	for i := range edges {
		e0 := edges[i][0]
		e1 := edges[i][1]
		if _, ok := g.adjacency[e0]; !ok {
			g.adjacency[e0] = make(map[string]struct{})
		}
		g.adjacency[e0][e1] = struct{}{}
	}
	return &g
}

type deque struct {
	q []string
}

func (d *deque) notEmpty() bool {
	return len(d.q) > 0
}

func (d *deque) append(element string) {
	d.q = append(d.q, element)
}

func (d *deque) popLeft() string {
	e, rest := d.q[0], d.q[1:]
	d.q = rest
	return e
}

type set map[string]struct{}

func (s *set) remove(e string) {
	delete(*s, e)
}

func (s *set) put(e string) {
	(*s)[e] = struct{}{}
}

func (s *set) contains(e string) bool {
	_, found := (*s)[e]
	return found
}

type predecessor map[string][]string
type predecessorEdge map[string]*string

func (pe *predecessorEdge) get(e string) *string {
	s, ok := (*pe)[e]
	if !ok || s == nil {
		return nil
	}
	return s
}

type recentEdge map[string]*Edge

func (re *recentEdge) contains(u string, v string) bool {
	a, ok := (*re)[u]
	if !ok || a == nil {
		return false
	}
	return a[0] == v || a[1] == v
}

type distance map[string]float64

func (d *distance) getOrInf(e string) float64 {
	v, ok := (*d)[e]
	if !ok {
		return math.Inf(0)
	}
	return v
}

func (d *distance) contains(e string) bool {
	_, ok := (*d)[e]
	return ok
}

type counter map[string]int64

func (c *counter) get(e string) int64 {
	i, ok := (*c)[e]
	if !ok {
		return 0
	}
	return i
}

var ErrNotFound = errors.New("negative cycle not found")

func bellmanFord(
	g *DirectedGraph,
	from string,
	weightFunc Weight,
	pred predecessor,
	heuristic bool) (string, error) {
	q := deque{}
	q.append(from)

	inQueue := make(set)
	inQueue.put(from)

	predEdge := make(predecessorEdge)
	predEdge[from] = nil

	dist := make(distance)
	dist[from] = 0

	recentUpdate := make(recentEdge)
	recentUpdate[from] = nil

	count := make(counter)

	n := int64(len(g.adjacency))

	for q.notEmpty() {
		u := q.popLeft()
		inQueue.remove(u)
		skip := false
		for i := range pred[u] {
			if inQueue.contains(pred[u][i]) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		distU := dist[u]
		for v := range g.adjacency[u] {
			distV := distU + weightFunc(u, v)
			if distV < dist.getOrInf(v) {
				if heuristic {
					if recentUpdate.contains(u, v) {
						pred[v] = append(pred[v], u)
						return v, nil
					}
					if s := predEdge.get(v); s != nil && *s == u {
						recentUpdate[v] = recentUpdate[u]
					} else {
						recentUpdate[v] = &Edge{u, v}
					}
				}
				if !inQueue.contains(v) {
					q.append(v)
					inQueue.put(v)
					countV := count.get(v) + 1
					if countV == n {
						return v, nil
					}
					count[v] = countV
				}
				dist[v] = distV
				pred[v] = []string{u}
				predEdge[v] = &u
			} else if dist.contains(v) && distV == dist[v] {
				pred[v] = append(pred[v], u)
			}
		}
	}
	return "", ErrNotFound
}

type stack []*stackElement

func (s *stack) notEmpty() bool {
	return len(*s) > 0
}

func (s *stack) last() *stackElement {
	return (*s)[len(*s)-1]
}

func (s *stack) pop() (stack, *stackElement) {
	s1, rest := (*s)[len(*s)-1], (*s)[:len(*s)-1]
	return rest, s1
}

type stackElement struct {
	v     string
	preds *preds
}

func newStackElement(v string, srcPreds []string) *stackElement {
	newPreds := make([]string, len(srcPreds))
	copy(newPreds, srcPreds)
	return &stackElement{v, &preds{newPreds}}
}

type preds struct {
	p []string
}

func (p *preds) contains(e string) bool {
	for i := range p.p {
		if p.p[i] == e {
			return true
		}
	}
	return false
}

func (p *preds) pop() string {
	p1, rest := p.p[len(p.p)-1], p.p[:len(p.p)-1]
	p.p = rest
	return p1
}

func (p *preds) notEmpty() bool {
	return len(p.p) > 0
}

type path []string

func (p *path) reverse() {
	pp := *p
	for i, j := 0, len(pp)-1; i < j; i, j = i+1, j-1 {
		pp[i], pp[j] = pp[j], pp[i]
	}
}

func (p *path) notEmpty() bool {
	return len(*p) > 0
}

func (p *path) pop() (path, string) {
	pp := *p
	p1, rest := pp[len(pp)-1], pp[:len(pp)-1]
	return rest, p1
}

func FindNegativeCycle(
	g *DirectedGraph,
	from string,
	weightFunc Weight,
) ([]string, error) {
	pred := make(predecessor)
	v, err := bellmanFord(g, from, weightFunc, pred, true)
	if err != nil {
		return []string{}, err
	}
	var negCycle path
	var stack stack
	stack = append(stack, newStackElement(v, pred[v]))
	seen := make(set)
	seen.put(v)
	for stack.notEmpty() {
		elm := stack.last()
		node, preds := elm.v, elm.preds
		if preds.contains(v) {
			negCycle = append(negCycle, node, v)
			negCycle.reverse()
			return negCycle, nil
		}
		if preds.notEmpty() {
			nbr := preds.pop()
			if !seen.contains(nbr) {
				stack = append(stack, newStackElement(nbr, pred[nbr]))
				negCycle = append(negCycle, node)
				seen.put(nbr)
			}
		} else {
			stack, _ = stack.pop()
			if negCycle.notEmpty() {
				negCycle, _ = negCycle.pop()
			} else {
				if adj, ok := g.adjacency[v]; ok {
					if _, ok = adj[v]; ok && weightFunc(v, v) < 0 {
						return []string{v, v}, nil
					}
				}
				return []string{}, ErrNotFound
			}
		}
	}
	return []string{}, ErrNotFound
}
