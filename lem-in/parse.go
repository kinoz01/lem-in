package lemin

import (
	"container/list"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func GetGraph() *Graph {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: program input_file")
		os.Exit(1)
	}

	graph, _, err := ReadFile(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return graph
}

// Reads the input file and constructs the adjacency list, start/end rooms, and number of ants.
func ReadFile(filePath string) (*Graph, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		file, err = os.Open("./lemin_test/audit/" + filePath)
		if err != nil {
			return nil, "", fmt.Errorf("can't open your input file")
		}
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("can't read your input file")
	}

	lines := strings.Split(strings.TrimSpace(string(fileBytes)), "\n")

	graph := &Graph{Rooms: make(map[string]*Node)}
	graph.Exits = list.New()
	var startFound, endFound bool

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if i == 0 {
			var parseErr error
			graph.Ants, parseErr = strconv.Atoi(line)
			if parseErr != nil || graph.Ants == 0 {
				return nil, "", fmt.Errorf("error reading ants number from your input file")
			}
			continue
		}
		if line == "##start" {
			if startFound {
				return nil, "", fmt.Errorf("can't have more than one start")
			}
			graph.Start, err = ParseStartEnd("Start", i, lines)
			if err != nil {
				return nil, "", err
			}
			startFound = true
			continue
		}
		if line == "##end" {
			if endFound {
				return nil, "", fmt.Errorf("can't have more than one end")
			}
			graph.End, err = ParseStartEnd("End", i, lines)
			if err != nil {
				return nil, "", err
			}
			endFound =true
			continue
		}
		if  strings.HasPrefix(line, "L") {
			return nil, "", fmt.Errorf("can't start a room name with L")
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "-")
		if len(parts) == 2 {
			from := parts[0]
			to := parts[1]
			if from == to {
				continue
			}
			// Add nodes and edges to the graph
			if graph.Rooms[from] == nil {
				graph.Rooms[from] = &Node{Edges: make(map[string]bool), Prev: "L"}
			}
			if graph.Rooms[to] == nil {
				graph.Rooms[to] = &Node{Edges: make(map[string]bool), Prev: "L"}
			}
			graph.Rooms[from].Edges[to] = true
			graph.Rooms[to].Edges[from] = true
		}
	}

	// Validate graph structure
	if graph.Start == graph.End || graph.Start == "" || graph.End == "" {
		return nil, "", fmt.Errorf("wrong start/end room")
	}
	if len(graph.Rooms) == 0 {
		return nil, "", fmt.Errorf("can't find linked rooms")
	}
	// Check if start/end rooms are linked to a nother node.
	if _, startExist := graph.Rooms[graph.Start]; !startExist {
		return nil, "", fmt.Errorf("start room isn't linked")
	}
	if _, endExist := graph.Rooms[graph.End]; !endExist {
		return nil, "", fmt.Errorf("end room isn't linked")
	}

	return graph, strings.TrimSpace(string(fileBytes)), nil
}

// Parses the start or end room from the input.
func ParseStartEnd(which string, i int, lines []string) (string, error) {
	if i == len(lines)-1 {
		return "", fmt.Errorf("%s room is missing", which)
	}
	roomDef := strings.Fields(lines[i+1])
	if len(roomDef) != 3 || roomDef == nil {
		return "", fmt.Errorf("%s room coordinates are not correctly formatted", which)
	}
	roomName := roomDef[0]
	if strings.HasPrefix(roomName, "#") || strings.HasPrefix(roomName, "L") || strings.Contains(roomName, " ") {
		return "", fmt.Errorf("%s room is missing", which)
	}
	return roomName, nil
}
