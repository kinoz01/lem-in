package lemin

import "fmt"

// PrintGraph prints all information about the Graph
func PrintGraph(graph *Graph) {
	fmt.Println("Graph:")
	fmt.Printf("Start Room: %s\n", graph.Start)
	fmt.Printf("End Room: %s\n", graph.End)
	fmt.Printf("Ants: %d\n", graph.Ants)
	fmt.Println("Rooms:")
	for roomName, room := range graph.Rooms {
		fmt.Printf("Room: %s\n", roomName)
		fmt.Printf("  Edges: %v\n", room.Edges)
		fmt.Printf("  Prev: %s\n", room.Prev)
		fmt.Printf("  EdgeIn: %s, EdgeOut: %s\n", room.EdgeIn, room.EdgeOut)
		fmt.Printf("  PriceIn: %d, PriceOut: %d\n", room.PriceIn, room.PriceOut)
		fmt.Printf("  CostIn: %d, CostOut: %d\n", room.CostIn, room.CostOut)
		fmt.Printf("  Split: %t\n", room.Split)
	}
	fmt.Println("Exits:")
	for e := graph.Exits.Front(); e != nil; e = e.Next() {
		fmt.Printf("  %v\n", e.Value)
	}
}

// PrintPaths prints all information about the Paths
func PrintPaths(paths *Paths) {
	fmt.Println("Paths:")
	fmt.Printf("Number of Paths: %d\n", paths.NumPaths)
	fmt.Printf("Total Steps: %d\n", paths.TotalSteps)
	fmt.Println("All Paths:")
	for i, path := range paths.AllPaths {
		fmt.Printf("  Path %d: ", i+1)
		for p := path.Front(); p != nil; p = p.Next() {
			fmt.Printf("%v -> ", p.Value)
		}
		fmt.Println("end")
	}
	fmt.Println("Ant Assignment:")
	for i, ants := range paths.Assignment {
		fmt.Printf("  Path %d: %d ants\n", i+1, ants)
	}
}

// PrintPriorityQueue prints all information about the PriorityQueue
func PrintPriorityQueue(pq PriorityQueue) {
	fmt.Println("Priority Queue:")
	for i, node := range pq {
		fmt.Printf("  Node %d: Room=%s, Cost=%d, Index=%d\n", i+1, node.Room, node.Cost, node.Index)
	}
}