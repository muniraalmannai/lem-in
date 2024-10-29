package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Function to open the file, and save each line of its contents as a slice of string
func readFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Using bufio.scanner to read the file line by line
	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// Define a room with its name, coordinates, flags, and connected rooms
type Room struct {
	name     string
	x, y     int
	isStart  bool
	isEnd    bool
	adjacent []string
}

type FlowGraph struct {
	capacity map[string]map[string]int // Capacity between rooms (room1 -> room2)
	flow     map[string]map[string]int // Flow between rooms
}

// Defining global variables
var numAnts int
var rooms = make(map[string]*Room)
var startRoomName, endRoomName string // Global start and end room names
var graph = FlowGraph{
	capacity: make(map[string]map[string]int),
	flow:     make(map[string]map[string]int),
}

// Parse the number of ants from the first line
func parseAnts(line string) error {
	ants := 0
	_, err := fmt.Sscanf(line, "%d", &ants) //Extract number of ants from the string line
	if err != nil {
		return fmt.Errorf("Invalid number of ants")
	}
	numAnts = ants // Assign the parsed number of ants to the global variable instead of var ants
	return nil
}

func parseData(lines []string) error {
	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])

		// Skip comments
		if strings.HasPrefix(line, "#") {
			if line == "##start" && i+1 < len(lines) {
				i++
				name := ""
				x, y := 0, 0
				fmt.Sscanf(lines[i], "%s %d %d", &name, &x, &y)
				rooms[name] = &Room{name: name, x: x, y: y, isStart: true}
				startRoomName = name

				// Initialize the room in the graph if not already present
				if _, exists := graph.capacity[name]; !exists {
					graph.capacity[name] = make(map[string]int)
					graph.flow[name] = make(map[string]int)
				}
				i++
				continue
			} else if line == "##end" && i+1 < len(lines) {
				i++
				name := ""
				x, y := 0, 0
				fmt.Sscanf(lines[i], "%s %d %d", &name, &x, &y)
				rooms[name] = &Room{name: name, x: x, y: y, isEnd: true}
				endRoomName = name

				// Initialize the room in the graph if not already present
				if _, exists := graph.capacity[name]; !exists {
					graph.capacity[name] = make(map[string]int)
					graph.flow[name] = make(map[string]int)
				}
				i++
				continue
			} else {
				i++
				continue
			}
		}

		// Parse rooms (find lines with coordinates)
		name := ""
		x, y := 0, 0
		if _, err := fmt.Sscanf(line, "%s %d %d", &name, &x, &y); err == nil {
			room := &Room{name: name, x: x, y: y}
			rooms[name] = room
			// Initialize adjacency lists for flow graph
			if _, exists := graph.capacity[name]; !exists {
				graph.capacity[name] = make(map[string]int)
				graph.flow[name] = make(map[string]int)
			}
			i++
			continue
		}

		// Parse tunnels (find links between rooms)
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) == 2 {
				room1, room2 := parts[0], parts[1]
				// Add bidirectional tunnels with capacity 1 (for the ants)
				rooms[room1].adjacent = append(rooms[room1].adjacent, room2)
				rooms[room2].adjacent = append(rooms[room2].adjacent, room1)
				graph.capacity[room1][room2] = 1
				graph.capacity[room2][room1] = 1
			} else {
				return fmt.Errorf("Invalid link format")
			}
		}
		i++
	}
	return nil
}

// Helper function to reconstruct the path from start to end
func reconstructPath(prev map[string]string, start string, end string) []string {
	path := []string{}
	for at := end; at != ""; at = prev[at] {
		path = append([]string{at}, path...)
		if at == start {
			break
		}
	}
	return path
}

// BFS to find a path from start to end room
func bfs(start string, end string) ([]string, error) {
	// Initialize BFS structures
	queue := []string{start}                // we start from "start" room and explore neighbours by level
	visited := map[string]bool{start: true} // keep track of visited rooms to not revisit them
	prev := map[string]string{}             // stores the room we came from, to reconstruct the path

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// if we reach the end, we found a path
		if current == end {
			return reconstructPath(prev, start, end), nil
		}

		for _, neighbour := range rooms[current].adjacent {
			if !visited[neighbour] && graph.capacity[current][neighbour] > graph.flow[current][neighbour] {
				visited[neighbour] = true
				prev[neighbour] = current
				queue = append(queue, neighbour)
			}
		}
	}
	return nil, fmt.Errorf("No augmenting path found")
}

// Edmonds-Karp algorithm to find maximum flow from start to end
func edmondsKarp(start string, end string) int {
	maxFlow := 0

	// find an augmenting path using bfs
	for {
		path, err := bfs(start, end)
		if err != nil {
			break // no more augmenting paths, end the loop
		}

		// determine the maximum flow capacity along the path (find the bottleneck)
		pathFlow := 1 // from the question

		// update flow along the paths
		for i := 0; i < len(path)-1; i++ {
			u, v := path[i], path[i+1]
			graph.flow[u][v] += pathFlow // forward direction
			graph.flow[v][u] -= pathFlow // backward direction
		}

		// the number of ants that made it to the end
		maxFlow += pathFlow
	}
	return maxFlow
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No input file provided")
		return
	}

	if len(os.Args) > 2 {
		fmt.Println("Error: Too many arguments provided")
		return
	}

	filename := os.Args[1]
	lines, err := readFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Print file contents
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println()

	// Parse number of ants
	err = parseAnts(lines[0])
	if err != nil {
		fmt.Println("Error parsing number of ants:", err)
	}

	// Parse rooms and tunnels, and identify start/end rooms
	err = parseData(lines[1:])
	if err != nil {
		fmt.Println("Error parsing data:", err)
		return
	}

	// Run the algorithm to find maxFlow
	maxFlow := edmondsKarp(startRoomName, endRoomName)

	fmt.Printf("Maximum number of ants that can reach the end: %d\n", maxFlow)
}
