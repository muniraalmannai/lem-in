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
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Parse rooms (find lines with coordinates)
		name := ""
		x, y := 0, 0
		if _, err := fmt.Sscanf(line, "%s %d %d", &name, &x, &y); err == nil {
			room := &Room{name: name, x: x, y: y}
			rooms[name] = room
			// Initialise adjacency lists for flow graph
			if _, exists := graph.capacity[name]; !exists {
				graph.capacity[name] = make(map[string]int)
				graph.flow[name] = make(map[string]int)
			}
			continue
		}

		// Parse tunnels (find links between rooms)
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) == 2 {
				room1, room2 := parts[0], parts[1]
				// Add 2 directional tunnels with capacity of 1 (for the ants)
				rooms[room1].adjacent = append(rooms[room1].adjacent, room2)
				rooms[room2].adjacent = append(rooms[room2].adjacent, room1)
				graph.capacity[room1][room2] = 1
				graph.capacity[room2][room1] = 1
			} else {
				return fmt.Errorf("Invalid link format")
			}

		}
	}
	return nil
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

	// // debugging: Print file contents
	// for _, line := range lines {
	// 	fmt.Println(line)
	// }
}
