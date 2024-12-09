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
	Edges             map[string]bool
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

func Run() {
	graph := GetGraph()
	paths := ComputePaths(graph)
	if paths == nil {
		fmt.Println("No paths found")
		os.Exit(1)
	}
	SimulateAnts(paths, graph.Ants)
}
