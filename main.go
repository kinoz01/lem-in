package main

import (
	"container/heap"
	"container/list"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const Infinity = 1000000

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

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: program input_file")
		os.Exit(1)
	}
	adjList, startRoom, endRoom, antCount, _, err := ReadData(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	graph := NewGraph()
	graph.Start = startRoom
	graph.End = endRoom
	graph.Ants = antCount
	graph.Exits = list.New()

	// Add edges to the nodes
	for from, neighbors := range adjList {
		graph.Rooms[from] = &Node{Prev: "L", Edges: make(map[string]byte)}
		for _, to := range neighbors {
			graph.Rooms[from].Edges[to] = 1
		}
	}
	fmt.Println(graph.Rooms)

	paths := ComputePaths(graph)
	if paths == nil {
		fmt.Println("No paths found")
		os.Exit(1)
	}
	SimulateAnts(paths, graph.Ants)
}

// ReadData reads the input file and constructs the adjacency list, start/end rooms, and number of ants.
func ReadData(filePath string) (adjList map[string][]string, start, end string, antCount int, input string, err error) {

	adjList = make(map[string][]string)

	file, err := os.Open(filePath)
	if err != nil {
		file, err = os.Open("./examples/" + filePath)
		if err != nil {
			return nil, "", "", 0, "", fmt.Errorf("can't open your input file")
		}
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, "", "", 0, "", fmt.Errorf("can't read your input file")
	}

	input = string(fileBytes)
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		//trim left
		if i == 0 {
			var parseErr error
			antCount, parseErr = strconv.Atoi(line)
			if parseErr != nil || antCount == 0 {
				return nil, "", "", 0, "", fmt.Errorf("error reading ants number from your input file")
			}
			continue
		}
		if line == "##start" {
			start, err = ParseStartOrEnd("Start", i, lines)
			if err != nil {
				return nil, "", "", 0, "", err
			}
			continue
		}
		if line == "##end" {
			end, err = ParseStartOrEnd("End", i, lines)
			if err != nil {
				return nil, "", "", 0, "", err
			}
			continue
		}
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "L") || strings.HasPrefix(line, " ") {
			continue
		}

		parts := strings.Split(line, "-")
		if len(parts) == 2 {
			from := parts[0]
			to := parts[1]
			adjList[from] = append(adjList[from], to)
			adjList[to] = append(adjList[to], from) // Add the reverse link
		}
	}

	if start == end || start == "" || end == "" {
		return nil, "", "", 0, "", fmt.Errorf("wrong start/end room")
	}

	return adjList, start, end, antCount, input, nil
}

// ParseStartOrEnd parses the start or end room from the input.
func ParseStartOrEnd(which string, idx int, lines []string) (string, error) {
	if idx == len(lines)-1 {
		return "", fmt.Errorf("%s room is missing", which)
	}
	roomDef := strings.Fields(lines[idx+1])
	if len(roomDef) != 3 || roomDef == nil {
		return "", fmt.Errorf("%s room coordinates are not correctly formatted", which)
	}
	roomName := roomDef[0]
	if strings.HasPrefix(roomName, "#") || strings.HasPrefix(roomName, "L") || strings.Contains(roomName, " ") {
		return "", fmt.Errorf("%s room is missing", which)
	}
	return roomName, nil
}

// NewGraph initializes a new graph.
func NewGraph() *Graph {
	return &Graph{Rooms: make(map[string]*Node)}
}

// Dijkstra's algorithm to find the shortest path.
func Dijkstra(graph *Graph) bool {
	pq := make(PriorityQueue, 0, 100)
	ResetGraph(graph)
	heap.Push(&pq, &PQNode{Cost: 0, Room: graph.Start})
	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*PQNode)
		v := current.Room
		for w := range graph.Rooms[v].Edges {
			RelaxEdge(graph, &pq, v, w)
		}
	}
	SetPrices(graph)
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
func RelaxEdge(graph *Graph, pq *PriorityQueue, v, w string) {
	nodeV := graph.Rooms[v]
	nodeW := graph.Rooms[w]
	if v == graph.End || w == graph.Start || nodeW.Prev == v {
		return
	}
	if nodeV.Prev == w && nodeV.CostIn < Infinity && (1+nodeW.CostOut > nodeV.CostIn+nodeV.PriceIn-nodeW.PriceOut) {
		nodeW.EdgeOut = v
		nodeW.CostOut = nodeV.CostIn - 1 + nodeV.PriceIn - nodeW.PriceOut
		heap.Push(pq, &PQNode{Cost: nodeW.CostOut, Room: w})
		RelaxHiddenEdge(graph, pq, w)
	} else if nodeV.Prev != w && nodeV.CostOut < Infinity && -1+nodeW.CostIn > nodeV.CostOut+nodeV.PriceOut-nodeW.PriceIn {
		nodeW.EdgeIn = v
		nodeW.CostIn = nodeV.CostOut + 1 + nodeV.PriceOut - nodeW.PriceIn
		heap.Push(pq, &PQNode{Cost: nodeW.CostIn, Room: w})
		RelaxHiddenEdge(graph, pq, w)
	}
}

// RelaxHiddenEdge further relaxes edges for nodes that have been split.
func RelaxHiddenEdge(graph *Graph, pq *PriorityQueue, w string) {
	node := graph.Rooms[w]
	if node.Split && node.CostIn > node.CostOut+node.PriceOut-node.PriceIn && w != graph.Start {
		node.EdgeIn = node.EdgeOut
		node.CostIn = node.CostOut + node.PriceOut - node.PriceIn
		if node.CostIn != node.CostOut {
			heap.Push(pq, &PQNode{Cost: node.CostIn, Room: w})
		}
	}
	if !node.Split && node.CostOut > node.CostIn+node.PriceIn-node.PriceOut && w != graph.End {
		node.EdgeOut = node.EdgeIn
		node.CostOut = node.CostIn + node.PriceIn - node.PriceOut
		if node.CostIn != node.CostOut {
			heap.Push(pq, &PQNode{Cost: node.CostOut, Room: w})
		}
	}
}

// SetPrices updates the node prices after Dijkstra's algorithm.
func SetPrices(graph *Graph) {
	for _, node := range graph.Rooms {
		node.PriceIn = node.CostIn
		node.PriceOut = node.CostOut
	}
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

// ComputePaths computes all possible paths using Suurballe's algorithm.
func ComputePaths(graph *Graph) *Paths {
	var bestPaths *Paths
	var newPaths *Paths
	if bestPaths = GetNextPaths(graph); bestPaths == nil {
		return nil
	}
	pathCount := 1
	for pathCount < graph.Ants {
		if newPaths = GetNextPaths(graph); newPaths == nil {
			break
		}
		if newPaths.TotalSteps < bestPaths.TotalSteps {
			bestPaths = newPaths
		}
		pathCount++
	}
	return bestPaths
}

// GetNextPaths finds the next set of paths.
func GetNextPaths(graph *Graph) *Paths {
	if !Dijkstra(graph) {
		return nil
	}
	CachePath(graph)
	return PathsFromGraph(graph)
}

// PathsFromGraph constructs the paths from the graph.
func PathsFromGraph(graph *Graph) *Paths {
	paths := new(Paths)
	paths.NumPaths = graph.Exits.Len()
	paths.AllPaths = make([]*list.List, paths.NumPaths)
	i := 0
	for link := graph.Exits.Front(); link != nil; link = link.Next() {
		p := UnrollPath(graph, link.Value.(string))
		paths.AllPaths[i] = p
		i++
	}
	sort.Slice(paths.AllPaths, func(i, j int) bool { return paths.AllPaths[i].Len() < paths.AllPaths[j].Len() })
	paths.TotalSteps = paths.calculateSteps(graph.Ants)
	return paths
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

// distributeAnts assigns ants to paths to minimize total steps.
func (paths *Paths) distributeAnts(antCount int) {
	paths.Assignment = make([]int, paths.NumPaths)
	l := len(paths.AllPaths) - 1
	longest := paths.pathLength(l)
	var sum int
	for i := 0; i < paths.NumPaths; i++ {
		sum += longest - paths.pathLength(i)
	}
	avgAnts := float32(antCount-sum) / float32(paths.NumPaths)
	rem := (avgAnts - float32(int(avgAnts))) * float32(paths.NumPaths)
	for i := 0; i < paths.NumPaths; i++ {
		paths.Assignment[i] = longest - paths.pathLength(i) + int(avgAnts)
		if rem > 0 {
			paths.Assignment[i]++
			rem--
		}
	}
}

// SimulateAnts simulates the movement of ants along the paths and prints the steps.
func SimulateAnts(paths *Paths, antCount int) {
	paths.distributeAnts(antCount) // Distribute ants into each path
	var lastAnt int
	antNum, activeAnt := 1, 1
	antPositions := make(map[int]*list.Element)
	for j := 0; j < paths.TotalSteps; j++ {
		for k := activeAnt; k <= lastAnt; k++ {
			if pos, ok := antPositions[k]; ok && pos != nil {
				fmt.Printf("L%d-%v ", k, pos.Value)
				antPositions[k] = pos.Next()
			} else {
				activeAnt = k + 1
				delete(antPositions, k)
			}
		}
		for i := 0; i < paths.NumPaths; i++ {
			if antNum > antCount {
				break
			}
			if paths.Assignment[i] <= 0 {
				continue
			} else {
				paths.Assignment[i]--
			}
			nextRoom := paths.AllPaths[i].Front().Next()
			if nextRoom != nil {
				fmt.Printf("L%d-%v ", antNum, nextRoom.Value)
				antPositions[antNum] = nextRoom.Next()
			}
			antNum++
		}
		if len(antPositions) > 0 {
			fmt.Println()
		}
		lastAnt = antNum - 1
	}
}

// Implementation of heap.Interface for PriorityQueue
func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Cost < pq[j].Cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PQNode)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // Avoid memory leak
	item.Index = -1 // For safety
	*pq = old[0 : n-1]
	return item
}
