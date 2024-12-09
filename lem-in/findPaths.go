package lemin

import (
	"container/heap"
	"container/list"
	"sort"
	"strings"
)

// ComputePaths computes all possible paths using Suurballe's algorithm.
func ComputePaths(graph *Graph) *Paths {

	var bestPaths, newPaths *Paths
	if bestPaths = GetNextPaths(graph); bestPaths == nil {
		return nil
	}

	minPathFound := 1
	for minPathFound < graph.Ants {
		if newPaths = GetNextPaths(graph); newPaths == nil {
			break
		}

		if newPaths.TotalSteps < bestPaths.TotalSteps {
			bestPaths = newPaths
		}
		minPathFound++
	}

	return bestPaths
}

// GetNextPaths finds the next set of paths.
func GetNextPaths(graph *Graph) *Paths {
	if !Dijkstra(graph) {
		return nil
	}
	SetPrices(graph)
	CachePath(graph)
	return PathsFromGraph(graph)
}

// Dijkstra's algorithm to find the shortest path.
func Dijkstra(graph *Graph) bool {
	pq := make(PriorityQueue, 0, 100)
	ResetGraph(graph)
	heap.Push(&pq, &PQNode{Cost: 0, Room: graph.Start})

	for pq.Len() > 0 {
		currentNode := heap.Pop(&pq).(*PQNode).Room

		for neighbor := range graph.Rooms[currentNode].Edges {
			RelaxEdge(graph, &pq, currentNode, neighbor)
		}
	}
	return graph.Rooms[graph.End].EdgeIn != "L"
}

// ResetGraph resets the graph costs and parents before running Dijkstra's algorithm.
func ResetGraph(graph *Graph) {
	for _, node := range graph.Rooms {
		node.EdgeIn = "L"
		node.EdgeOut = "L"
		node.CostIn = Infinity
		node.CostOut = Infinity
	}
	graph.Rooms[graph.Start].CostIn = 0
	graph.Rooms[graph.Start].CostOut = 0
}

// RelaxEdge relaxes the edges during Dijkstra's algorithm.
func RelaxEdge(graph *Graph, pq *PriorityQueue, current, next string) {
	currentNode := graph.Rooms[current]
	nextNode := graph.Rooms[next]

	if current == graph.End || next == graph.Start || nextNode.Prev == current {
		return
	}

	if currentNode.Prev == next && currentNode.CostIn+currentNode.PriceIn < nextNode.CostOut+nextNode.PriceOut+1 {
		nextNode.EdgeOut = current
		nextNode.CostOut = currentNode.CostIn - 1 + currentNode.PriceIn - nextNode.PriceOut
		heap.Push(pq, &PQNode{Cost: nextNode.CostOut, Room: next})
		RelaxHiddenEdge(graph, pq, next)
	} else if currentNode.Prev != next && currentNode.CostOut+currentNode.PriceOut+1 < nextNode.CostIn+nextNode.PriceIn {
		nextNode.EdgeIn = current
		nextNode.CostIn = currentNode.CostOut + 1 + currentNode.PriceOut - nextNode.PriceIn
		heap.Push(pq, &PQNode{Cost: nextNode.CostIn, Room: next})
		RelaxHiddenEdge(graph, pq, next)
	}
}

// RelaxHiddenEdge further relaxes edges for nodes that have been split.
func RelaxHiddenEdge(graph *Graph, pq *PriorityQueue, nextNode string) {
	node := graph.Rooms[nextNode]
	if node.Split && node.CostIn+node.PriceIn > node.CostOut+node.PriceOut && nextNode != graph.Start {
		node.EdgeIn = node.EdgeOut
		node.CostIn = node.CostOut + node.PriceOut - node.PriceIn
		if node.CostIn != node.CostOut {
			heap.Push(pq, &PQNode{Cost: node.CostIn, Room: nextNode})
		}
	}
	if !node.Split && node.CostOut+node.PriceOut > node.CostIn+node.PriceIn && nextNode != graph.End {
		node.EdgeOut = node.EdgeIn
		node.CostOut = node.CostIn + node.PriceIn - node.PriceOut
		if node.CostIn != node.CostOut {
			heap.Push(pq, &PQNode{Cost: node.CostOut, Room: nextNode})
		}
	}
}

// PathsFromGraph constructs the paths from the graph.
func PathsFromGraph(graph *Graph) *Paths {
	paths := new(Paths)
	uniquePaths := make(map[string]*list.List)
	for link := graph.Exits.Front(); link != nil; link = link.Next() {
		p := UnrollPath(graph, link.Value.(string))
		pathStr := PathToString(p)
		uniquePaths[pathStr] = p
	}
	// Convert the map to a slice
	paths.AllPaths = make([]*list.List, 0, len(uniquePaths))
	for _, p := range uniquePaths {
		paths.AllPaths = append(paths.AllPaths, p)
	}
	paths.NumPaths = len(paths.AllPaths)
	sort.Slice(paths.AllPaths, func(i, j int) bool { return paths.AllPaths[i].Len() < paths.AllPaths[j].Len() })
	paths.TotalSteps = paths.calculateSteps(graph.Ants)
	return paths
}

func PathToString(path *list.List) string {
	var nodes []string
	for e := path.Front(); e != nil; e = e.Next() {
		nodes = append(nodes, e.Value.(string))
	}
	return strings.Join(nodes, "->")
}

// pathLength returns the length of the i-th path.
func (paths *Paths) pathLength(i int) int {
	return paths.AllPaths[i].Len()
}

// calculateSteps calculates the total steps required for all ants to reach the end.
func (paths *Paths) calculateSteps(antCount int) int {
	l := len(paths.AllPaths) - 1
	shortest := paths.pathLength(0)
	longest := paths.pathLength(l)
	var sum int
	for i := 0; i < paths.NumPaths; i++ {
		sum += longest - paths.pathLength(i)
	}
	antsPerPath := longest - shortest + (antCount-sum)/paths.NumPaths
	if (antCount-sum)%paths.NumPaths > 0 {
		antsPerPath++
	}
	return shortest + antsPerPath - 1
}

// UnrollPath reconstructs a path from the end node to the start node.
func UnrollPath(graph *Graph, v string) *list.List {
	path := list.New()
	path.PushFront(graph.End)
	for v != graph.Start {
		path.PushFront(v)
		v = graph.Rooms[v].Prev
	}
	path.PushFront(graph.Start)
	return path
}

// CachePath caches the path found by Dijkstra's algorithm.
func CachePath(graph *Graph) {
	var unsplit bool
	w := graph.End
	v := graph.Rooms[w].EdgeIn
	graph.Exits.PushBack(v)
	for w != graph.Start {
		if graph.Rooms[v].Prev == w {
			if unsplit {
				UnsplitNode(graph, w)
			}
			unsplit = true
			w, v = v, graph.Rooms[v].EdgeIn
		} else {
			graph.Rooms[w].Prev = v
			SplitNode(graph, w)
			unsplit = false
			w, v = v, graph.Rooms[v].EdgeOut
		}
	}
}

// UnsplitNode resets a split node.
func UnsplitNode(graph *Graph, v string) {
	graph.Rooms[v].Split = false
	graph.Rooms[v].Prev = "L"
}

// SplitNode marks a node as split to prevent edge reuse.
func SplitNode(graph *Graph, v string) {
	if v != graph.Start && v != graph.End {
		graph.Rooms[v].Split = true
	}
}

// SetPrices updates the node prices after Dijkstra's algorithm.
func SetPrices(graph *Graph) {
	for _, node := range graph.Rooms {
		node.PriceIn = node.CostIn
		node.PriceOut = node.CostOut
	}
}
