package lemin

import (
	"container/list"
	"fmt"
	"os"
)

const Infinity = 1 << 60

// Graph represents the network of rooms and connections.
type Graph struct {
	Rooms      map[string]*Node
	Exits      *list.List
	Start, End string
	Ants       int
}

// Node represents a room in the graph.
type Node struct {
	Edges             map[string]byte
	Prev              string
	EdgeIn, EdgeOut   string
	PriceIn, PriceOut int
	CostIn, CostOut   int
	Split             bool
}

// Paths holds information about possible paths and ant assignments.
type Paths struct {
	NumPaths, TotalSteps int
	AllPaths             []*list.List
	Assignment           []int // Number of ants assigned to each path
}

// PQNode is a node in the priority queue used in Dijkstra's algorithm.
type PQNode struct {
	Cost  int
	Index int
	Room  string
}

// PriorityQueue implements heap.Interface and holds PQNodes.
type PriorityQueue []*PQNode

func Run() {
	graph := GetGraph()
	PrintGraph(graph)
	fmt.Println("++++++++++++++")
	paths := ComputePaths(graph)
	if paths == nil {
		fmt.Println("No paths found")
		os.Exit(1)
	}

	SimulateAnts(paths, graph.Ants)
}

/* for i, l := range paths.AllPaths {
	fmt.Printf("List %d:\n", i+1)
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
} */
