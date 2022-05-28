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

type stack struct {
	s []*stackElement
}

func (s *stack) notEmpty() bool {
	return len(s.s) > 0
}

func (s *stack) last() *stackElement {
	return s.s[len(s.s)-1]
}

func (s *stack) append(e *stackElement) {
	s.s = append(s.s, e)
}

func (s *stack) pop() *stackElement {
	s1 := s.s[len(s.s)-1]
	rest := s.s[:len(s.s)-1]
	s.s = rest
	return s1
}

type stackElement struct {
	node  string
	preds *preds
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
	p1 := p.p[len(p.p)-1]
	rest := p.p[:len(p.p)-1]
	p.p = rest
	return p1
}

func (p *preds) notEmpty() bool {
	return len(p.p) > 0
}

type path struct {
	p []string
}

func (p *path) reverse() {
	for i, j := 0, len(p.p)-1; i < j; i, j = i+1, j-1 {
		p.p[i], p.p[j] = p.p[j], p.p[i]
	}
}

func (p *path) notEmpty() bool {
	return len(p.p) > 0
}

func (p *path) append(r string) {
	p.p = append(p.p, r)
}

func (p *path) pop() string {
	p1 := p.p[len(p.p)-1]
	rest := p.p[:len(p.p)-1]
	p.p = rest
	return p1
}

func (p *path) copy() []string {
	pp := make([]string, len(p.p))
	copy(pp, p.p)
	return pp
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
	stack.append(&stackElement{v, &preds{pred[v]}})
	seen := make(set)
	seen.put(v)
	for stack.notEmpty() {
		elm := stack.last()
		// node, preds := elm.v, elm.preds
		if elm.preds.contains(v) {
			negCycle.append(elm.node)
			negCycle.append(v)
			negCycle.reverse()
			return negCycle.copy(), nil
		}
		if elm.preds.notEmpty() {
			nbr := elm.preds.pop()
			if !seen.contains(nbr) {
				stack.append(&stackElement{nbr, &preds{pred[nbr]}})
				negCycle.append(elm.node)
				seen.put(nbr)
			}
		} else {
			stack.pop()
			if negCycle.notEmpty() {
				negCycle.pop()
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

func StartFrom(path []string, start string) []string {
	if path[len(path)-1] != path[0] {
		path = append(path, path[0])
	}
	startIndex := -1
	for i := range path {
		if path[i] == start {
			startIndex = i
			break
		}
	}
	newPath := make([]string, len(path))
	for i := startIndex; i < len(path)-1; i++ {
		newPath[i-startIndex] = path[i]
	}
	for i := 0; i < startIndex+1; i++ {
		newPath[i+len(path)-startIndex-1] = path[i]
	}
	return newPath
}
