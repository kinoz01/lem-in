package lemin

import "fmt"

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

func PrintPaths(paths *Paths) {
	fmt.Println("Paths:")
	fmt.Printf("Number of Paths: %d\n", paths.NumPaths)
	fmt.Printf("Total Steps: %d\n", paths.TotalSteps)
	fmt.Println("All Paths:")
	for i, path := range paths.AllPaths {
		fmt.Printf("  Path %d: ", i+1)
		for p := path.Front(); p != nil; p = p.Next() {
			fmt.Printf("%v ", p.Value)
			if p.Next() != nil {
				fmt.Print("-> ")
			}
		}
		fmt.Println()
	}
	fmt.Println("Ant Assignment:")
	for i, ants := range paths.Assignment {
		fmt.Printf("  Path %d: %d ants\n", i+1, ants)
	}
}

func PrintPriorityQueue(pq PriorityQueue) {
	fmt.Println("Priority Queue:")
	for i, node := range pq {
		fmt.Printf("  Node %d: Room=%s, Cost=%d, Index=%d\n", i+1, node.Room, node.Cost, node.Index)
	}
}

func PrintNode(node *Node, name string) {
	fmt.Printf("Node: %s\n", name)
	fmt.Println("Edges:", node.Edges)
	fmt.Printf("Prev: %s\n", node.Prev)
	fmt.Printf("EdgeIn: %s, EdgeOut: %s\n", node.EdgeIn, node.EdgeOut)
	fmt.Printf("PriceIn: %d, PriceOut: %d\n", node.PriceIn, node.PriceOut)
	fmt.Printf("CostIn: %d, CostOut: %d\n", node.CostIn, node.CostOut)
	fmt.Printf("Split: %v\n", node.Split)
	fmt.Println("nodenodenodenodenode")
}
