package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var rooms = make(map[string]map[string]int) // Map to store connections between rooms
var startRoom, endRoom string               // Start and end room identifiers
var numAnts int                             // Number of ants to be processed

// Parses the input file and populates rooms, startRoom, endRoom, and numAnts
func parseFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if lineCount == 0 {
			// Parse the number of ants
			fmt.Sscanf(line, "%d", &numAnts)
			lineCount++
			continue
		}

		// Handle room definitions and links
		if line == "##start" {
			// Parse the start room
			scanner.Scan()
			startRoom = strings.Fields(scanner.Text())[0]
		} else if line == "##end" {
			// Parse the end room
			scanner.Scan()
			endRoom = strings.Fields(scanner.Text())[0]
		} else if strings.Contains(line, "-") {
			// Parse room connections (links)
			parts := strings.Split(line, "-")
			room1, room2 := parts[0], parts[1]

			// Initialize room maps if not present
			if rooms[room1] == nil {
				rooms[room1] = make(map[string]int)
			}
			if rooms[room2] == nil {
				rooms[room2] = make(map[string]int)
			}

			// Set connections between rooms
			rooms[room1][room2] = 1
			rooms[room2][room1] = 1
		}
	}

	return scanner.Err() // Return any scanning errors encountered
}

// Uses DFS to find all paths from start to end without revisiting rooms
func findAllPaths(start, end string, visited map[string]bool) [][]string {
	if start == end {
		return [][]string{{end}}
	}

	if visited[start] {
		return nil
	}

	visited[start] = true
	var allPaths [][]string

	// Explore each neighboring room
	for neighbor := range rooms[start] {
		if rooms[start][neighbor] > 0 {
			// Create a copy of the visited map
			newVisited := make(map[string]bool)
			for k, v := range visited {
				newVisited[k] = v
			}

			// Recursively find paths from the neighbor to the end room
			paths := findAllPaths(neighbor, end, newVisited)
			for _, path := range paths {
				newPath := append([]string{start}, path...)
				allPaths = append(allPaths, newPath)
			}
		}
	}

	visited[start] = false
	return allPaths
}

// Sorts paths by length (from shortest to longest)
func findOptimalPaths(paths [][]string) [][]string {
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})
	return paths
}

// Assigns ants to paths based on path length and current occupancy to optimize flow
func allocateAntsToPaths(paths [][]string) [][]int {
	antAssignments := make([][]int, len(paths))
	pathOccupancy := make([]int, len(paths))

	for ant := 0; ant < numAnts; ant++ {
		bestPathIndex := 0
		bestScore := len(paths[0]) + pathOccupancy[0]

		for i := 1; i < len(paths); i++ {
			score := len(paths[i]) + pathOccupancy[i]
			if score < bestScore {
				bestScore = score
				bestPathIndex = i
			}
		}

		// Assign ant to the best path
		antAssignments[bestPathIndex] = append(antAssignments[bestPathIndex], ant)
		pathOccupancy[bestPathIndex]++
	}

	return antAssignments
}

// Simulates and prints ant movements along their assigned paths
func simulateAntMovements(paths [][]string, antAssignments [][]int) {
	antPositions := make([]int, numAnts)  // Tracks each ant's position along its assigned path
	finishedAnts := make([]bool, numAnts) // Tracks whether each ant has reached the end
	totalFinishedAnts := 0

	for i := range antPositions {
		antPositions[i] = -1 // Initialize all ants at the start
	}

	for totalFinishedAnts < numAnts {
		var moves []string

		// Process each path
		for pathIndex, path := range paths {
			for _, ant := range antAssignments[pathIndex] {
				if finishedAnts[ant] {
					continue // Skip ants that have already reached the end
				}

				// Move the ant along the path
				if antPositions[ant] < len(path)-1 {
					antPositions[ant]++
					if antPositions[ant] > 0 {
						moves = append(moves, fmt.Sprintf("L%d-%s", ant+1, path[antPositions[ant]]))
					}

					// Mark as finished if the ant has reached the end
					if antPositions[ant] == len(path)-1 {
						finishedAnts[ant] = true
						totalFinishedAnts++
					}
				}
			}
		}

		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		} else {
			break
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No input file provided")
		return
	}

	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Println(string(content))
	fmt.Println()

	err = parseFile(filename)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	visited := make(map[string]bool)
	allPaths := findAllPaths(startRoom, endRoom, visited)

	if len(allPaths) == 0 {
		fmt.Println("Error: No paths found")
		return
	}

	optimalPaths := findOptimalPaths(allPaths)
	antAssignments := allocateAntsToPaths(optimalPaths)
	simulateAntMovements(optimalPaths, antAssignments)
}
