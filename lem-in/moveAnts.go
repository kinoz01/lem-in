package lemin

import (
	"container/list"
	"fmt"
)

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
		if len(antPositions) > 0 {
			fmt.Println()
		}
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

		lastAnt = antNum - 1
	}
}
