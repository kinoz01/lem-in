package main

import (
	"container/heap"
	"fmt"
	"math"
)

// Node represents a vertex in the graph.
type Node struct {
	Edges map[string]float64 // Neighbors with weights
	Prev  string             // Previous node in the path
	Cost  float64            // Cost to reach this node
}

// Graph represents the weighted graph.
type Graph struct {
	Rooms map[string]*Node // All nodes in the graph
	Start string           // Start node
	End   string           // End node
}

// PriorityQueue implementation
type PQNode struct {
	Cost float64
	Room string
}

type PriorityQueue []*PQNode

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Cost < pq[j].Cost
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*PQNode))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}


// Resets all nodes before running Dijkstra's algorithm.
func ResetGraph(graph *Graph) {
	for _, node := range graph.Rooms {
		node.Prev = ""
		node.Cost = math.Inf(1)
	}
	graph.Rooms[graph.Start].Cost = 0
}

// Relaxes an edge during Dijkstra's algorithm.
func RelaxEdge(graph *Graph, pq *PriorityQueue, v, w string) {
	nodeV := graph.Rooms[v]
	nodeW := graph.Rooms[w]
	weight := nodeV.Edges[w]

	if nodeV.Cost+weight < nodeW.Cost {
		nodeW.Cost = nodeV.Cost + weight
		nodeW.Prev = v
		heap.Push(pq, &PQNode{Cost: nodeW.Cost, Room: w})
	}
}

// Dijkstra runs the shortest path algorithm on the graph.
func Dijkstra(graph *Graph) []string {
	pq := &PriorityQueue{}

	ResetGraph(graph)
	heap.Push(pq, &PQNode{Cost: 0, Room: graph.Start})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*PQNode)
		v := current.Room
		for w := range graph.Rooms[v].Edges {
			RelaxEdge(graph, pq, v, w)
		}
	}

	// Reconstruct the shortest path
	path := []string{}
	current := graph.End
	for current != "" {
		path = append([]string{current}, path...)
		current = graph.Rooms[current].Prev
	}

	return path
}

func main() {
	graph := &Graph{
		Rooms: make(map[string]*Node),
		Start: "A",
		End:   "F",
	}

	// Add nodes
	graph.Rooms["A"] = &Node{Edges: make(map[string]float64)}
	graph.Rooms["B"] = &Node{Edges: make(map[string]float64)}
	graph.Rooms["C"] = &Node{Edges: make(map[string]float64)}
	graph.Rooms["D"] = &Node{Edges: make(map[string]float64)}
	graph.Rooms["E"] = &Node{Edges: make(map[string]float64)}
	graph.Rooms["F"] = &Node{Edges: make(map[string]float64)}

	// Add edges	
	graph.Rooms["A"].Edges["C"] = 35
	graph.Rooms["A"].Edges["D"] = 40
	graph.Rooms["B"].Edges["D"] = 20
	graph.Rooms["B"].Edges["E"] = 25
	graph.Rooms["C"].Edges["F"] = 30
	graph.Rooms["D"].Edges["F"] = 20
	graph.Rooms["D"].Edges["E"] = 45
	graph.Rooms["E"].Edges["F"] = 25
	graph.Rooms["A"].Edges["B"] = 5
	graph.Rooms["E"].Edges["C"] = 30

	// Make the graph bidirectional
	AddReversedEdges(graph)

	// Run Dijkstra's algorithm
	shortestPath := Dijkstra(graph)
	fmt.Println("Shortest Path:", shortestPath)
}

func PrintPriorityQueue(pq *PriorityQueue) {
	fmt.Println("PriorityQueue contents:")
	for _, item := range *pq {
		fmt.Printf("Room: %s, Cost: %.2f\n", item.Room, item.Cost)
	}
}

func AddReversedEdges(graph *Graph) {
	for nodeName, node := range graph.Rooms {
		for neighbor, weight := range node.Edges {
			// Check if the reverse edge already exists; if not, add it
			if _, exists := graph.Rooms[neighbor].Edges[nodeName]; !exists {
				graph.Rooms[neighbor].Edges[nodeName] = weight
			}
		}
	}
}
