package bf

func AllSimplePaths(g *DirectedGraph, from, to string, cutoff int) <-chan []string {
	ch := make(chan []string, 1)
	go func() {
		visited := edgeVisited{}
		visited.add(from)
		stack := []*edgeIterator{newEdgeIterator(g, from)}
		for len(stack) > 0 {
			children := stack[len(stack)-1]
			child := children.next()
			if child == nil {
				stack = stack[:len(stack)-1]
				visited.pop()
			} else if visited.length() < cutoff {
				if visited.contains(*child) {
					continue
				}
				if *child == to {
					var tmp []string
					tmp = append(tmp, visited.visitedOrders...)
					tmp = append(tmp, *child)
					ch <- tmp
				}
				visited.add(*child)
				if *child != to {
					stack = append(stack, newEdgeIterator(g, *child))
				} else {
					visited.pop()
				}
			} else {
				for i := range children.items {
					childItem := children.items[i]
					if visited.contains(childItem) {
						continue
					}
					if childItem == to {
						var tmp []string
						tmp = append(tmp, visited.visitedOrders...)
						tmp = append(tmp, childItem)
						ch <- tmp
					}
				}
				stack = stack[:len(stack)-1]
				visited.pop()
			}
		}
		close(ch)
	}()
	return ch
}

type edgeIterator struct {
	index int
	items []string
}

func newEdgeIterator(g *DirectedGraph, from string) *edgeIterator {
	ei := edgeIterator{}
	for k := range g.adjacency[from] {
		ei.items = append(ei.items, k)
	}
	return &ei
}

func (ei *edgeIterator) next() *string {
	if ei.index < len(ei.items) {
		i := ei.index
		ei.index += 1
		return &ei.items[i]
	}
	return nil
}

type edgeVisited struct {
	visitedOrders []string
	visitedMap    map[string]struct{}
}

func (ev *edgeVisited) add(e string) {
	if ev.visitedMap == nil {
		ev.visitedMap = make(map[string]struct{})
	}
	ev.visitedMap[e] = struct{}{}
	ev.visitedOrders = append(ev.visitedOrders, e)
}

func (ev *edgeVisited) contains(e string) bool {
	if ev.visitedMap == nil {
		return false
	}
	_, ok := ev.visitedMap[e]
	return ok
}

func (ev *edgeVisited) pop() {
	if ev.visitedMap == nil {
		return
	}
	e := ev.visitedOrders[len(ev.visitedOrders)-1]
	ev.visitedOrders = ev.visitedOrders[:len(ev.visitedOrders)-1]
	delete(ev.visitedMap, e)
}

func (ev *edgeVisited) length() int {
	return len(ev.visitedOrders)
}
