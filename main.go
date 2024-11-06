package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Farm struct represents the configuration of the ant farm.
type Farm struct {
	AntCount      int
	Start         string
	End           string
	AdjacencyList map[string][]string
	Paths         [][]string
}

// Limits to be set by user in main function
type PathLimits struct {
	MaxRoomsPerPath       int
	MaxTotalPaths         int
	MaxPathCombination    int
	MaxDirectPaths        int
	MaxPathsInCombination int
}

// loadFarm loads the configuration of the farm into Farm struct.
func loadFarm(filename string) (*Farm, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Scan the file and appends each line to 'lines'
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for empty file first (prevent panics)
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty input file")
	}

	// Extract ant count
	antCount, err := strconv.Atoi(lines[0])
	if err != nil || antCount <= 0 {
		return nil, fmt.Errorf("invalid number of ants: must be a number greater than 0")
	}

	var start, end string
	var links []string

	// Loop through each line to find the start, end, and room connections
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		switch line {
		case "##start":
			// Boundary (prevent panics)
			if i+1 >= len(lines) {
				return nil, fmt.Errorf("invalid file format: missing room definition after start command")
			}
			start = strings.Fields(lines[i+1])[0]
			i++ // The next line is already processed
		case "##end":
			// Boundary (prevent panics)
			if i+1 >= len(lines) {
				return nil, fmt.Errorf("invalid file format: missing room definition after end command")
			}
			end = strings.Fields(lines[i+1])[0]
			i++
		default:
			// If the line contains a dash, append into links slice
			if strings.Contains(line, "-") {
				links = append(links, line)
			}
		}
	}

	// Error checks before building adjacency list
	if start == "" {
		return nil, fmt.Errorf("no start room defined")
	}
	if end == "" {
		return nil, fmt.Errorf("no end room defined")
	}
	if len(links) == 0 {
		return nil, fmt.Errorf("no room connections found")
	}

	// Build the adjacency list from the links
	adjacencyList, err := buildAdjacencyList(links)
	if err != nil {
		return nil, err
	}
	return &Farm{AntCount: antCount, Start: start, End: end, AdjacencyList: adjacencyList}, nil
}

// buildAdjacencyList creates a map from links[]. Each link shows a direct connection between two rooms.
func buildAdjacencyList(edges []string) (map[string][]string, error) {
	adjList := make(map[string][]string)
	for _, edge := range edges {
		parts := strings.Split(edge, "-")
		if len(parts) != 2 || parts[0] == parts[1] {
			return nil, fmt.Errorf("invalid link between rooms: %s", edge)
		}
		adjList[parts[0]] = append(adjList[parts[0]], parts[1]) // Connect room A to room B
		adjList[parts[1]] = append(adjList[parts[1]], parts[0]) // Connect room B to room A
	}
	return adjList, nil
}

// findAllPaths performs Depth-First Search (DFS) to find all paths from start to end in the adjacency list.
func findAllPaths(adjacencyList map[string][]string, start, end string, limits PathLimits) [][]string {
	var paths [][]string

	var dfs func(current string, path []string)
	dfs = func(current string, path []string) {
		if len(path) > limits.MaxRoomsPerPath {
			return
		}

		path = append(path, current)
		if current == end {
			if len(paths) >= limits.MaxTotalPaths {
				return
			}
			paths = append(paths, append([]string{}, path...))
			return
		}

		for _, neighbor := range adjacencyList[current] {
			if !contains(path, neighbor) {
				dfs(neighbor, path)
			}
		}
	}
	dfs(start, []string{})

	if len(paths) == 0 {
		fmt.Println("Error: No valid paths found within length limit of", limits.MaxRoomsPerPath, "rooms")
		return paths
	}

	return paths
}

// findAllPaths helper, checks if a specific room is already in the current path to avoid loops.
func contains(path []string, room string) bool {
	for _, v := range path {
		if v == room {
			return true
		}
	}
	return false
}

// findOptimalPathCombination calculates the optimal combination of paths in the minimum number of turns.
func findOptimalPathCombination(paths [][]string, antCount int, limits PathLimits) ([][]string, []int) {
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	if len(paths) > limits.MaxPathCombination {
		fmt.Printf("Notice: Limiting path combinations to %d most efficient paths\n", limits.MaxPathCombination)
		paths = paths[:limits.MaxPathCombination]
	}

	if containsDirectPath(paths) {
		var optimalPaths [][]string
		for _, path := range paths {
			if isValidCombination(append(optimalPaths, path)) {
				optimalPaths = append(optimalPaths, path)
			}
			if len(optimalPaths) >= limits.MaxDirectPaths {
				break
			}
		}

		// Calculate the length of each path for turn calculations and track the distribution of ants across paths.
		pathLengths := make([]int, len(optimalPaths))
		for i := range optimalPaths {
			pathLengths[i] = len(optimalPaths[i]) - 1
		}
		antQueue := make([]int, len(optimalPaths))
		remaining := antCount

		// Distribute ants across paths, prioritizing the shortest path (in terms of turns) each time
		for remaining > 0 {
			shortestIdx := 0
			minTurns := pathLengths[0] + antQueue[0]

			// Find the path that results in the fewest turns with current distribution
			for i := 1; i < len(optimalPaths); i++ {
				turns := pathLengths[i] + antQueue[i]
				if turns < minTurns {
					minTurns = turns
					shortestIdx = i
				}
			}

			// Assign one ant to the selected shortest path
			antQueue[shortestIdx]++
			remaining-- // Decrease the remaining ants to assign
		}

		return optimalPaths, antQueue
	}

	// If no direct path is available, explore combinations to find the minimum turn configuration
	var optimalPaths [][]string
	var optimalAntQueue []int
	minTurns := int(^uint(0) >> 1)

	// Generate all valid path combinations and find the one with minimum total turns
	for i := 1; i <= len(paths); i++ {
		combinations := generateCombinations(paths, i, limits)
		// Combination checker:
		for _, combo := range combinations {
			if isValidCombination(combo) {
				// Calculate ant distribution and total turns for this combination
				antQueue, turns := distributeAnts(combo, antCount)
				// Update optimal configuration if this combination has fewer turns
				if turns < minTurns {
					minTurns = turns
					optimalPaths = combo
					optimalAntQueue = antQueue
				}
			}
		}
	}
	return optimalPaths, optimalAntQueue
}

// findOptimalPathCombination helper, checks if there is a direct path from start to end in the list of paths.
func containsDirectPath(paths [][]string) bool {
	for _, path := range paths {
		if len(path) == 2 {
			return true
		}
	}
	return false
}

// findOptimalPathCombination helper, checks that paths do not overlap.
func isValidCombination(paths [][]string) bool {
	roomSet := make(map[string]bool)
	for _, path := range paths {
		for _, room := range path[1 : len(path)-1] { // Ignore start and end rooms
			if roomSet[room] {
				return false
			}
			roomSet[room] = true
		}
	}
	return true
}

// generateCombinations generates all combinations of paths with a specified length.
func generateCombinations(paths [][]string, length int, limits PathLimits) [][][]string {
	if length > limits.MaxPathsInCombination {
		fmt.Println("Error: Exceeded maximum allowed path combination length")
		return nil
	}

	var combinations [][][]string
	comb := make([][]string, length)

	var generate func(int, int)
	generate = func(start, depth int) {
		if depth == length {
			combCopy := make([][]string, len(comb))
			copy(combCopy, comb)
			combinations = append(combinations, combCopy)
			return
		}
		for i := start; i < len(paths); i++ {
			comb[depth] = paths[i]
			generate(i+1, depth+1)
		}
	}
	generate(0, 0)
	return combinations
}

// distributeAnts assigns ants to paths, preferring shorter paths.
func distributeAnts(paths [][]string, antCount int) ([]int, int) {
	antQueue := make([]int, len(paths)) // Queue to track ants per path
	pathLengths := make([]int, len(paths))
	for i := range paths {
		pathLengths[i] = len(paths[i]) - 1 // Calculate length of each path
	}

	remaining := antCount
	for remaining > 0 {
		shortestIdx := findShortestPath(pathLengths, antQueue)
		antQueue[shortestIdx]++
		remaining--
	}

	return antQueue, maxPathTurns(pathLengths, antQueue)
}

// findShortestPath finds the path with the shortest distance plus ants assigned to it.
func findShortestPath(pathLengths, antQueue []int) int {
	minIndex := 0
	minScore := pathLengths[0] + antQueue[0]
	for i := 1; i < len(pathLengths); i++ {
		score := pathLengths[i] + antQueue[i]
		if score < minScore {
			minIndex = i
			minScore = score
		}
	}
	return minIndex
}

// maxPathTurns calculates the maximum number of turns needed for any ant to reach the end.
func maxPathTurns(pathLengths, antQueue []int) int {
	maxTurns := 0
	for i := range pathLengths {
		turns := pathLengths[i] + antQueue[i] - 1
		if turns > maxTurns {
			maxTurns = turns
		}
	}
	return maxTurns
}

// simulateAntMovement simulates the movement of each ant along the assigned paths.
func simulateAntMovement(paths [][]string, antQueue []int, antCount int) {
	type Ant struct {
		id       int
		pathIdx  int
		position int
		done     bool
	}

	ants := make([]Ant, antCount) // Create a slice for each ant
	antId, pathIdx := 1, 0

	// Initialize ants and assign them to paths
	for antId <= antCount {
		if antQueue[pathIdx] > 0 {
			ants[antId-1] = Ant{id: antId, pathIdx: pathIdx, position: 0}
			antQueue[pathIdx]--
			antId++
		}
		pathIdx = (pathIdx + 1) % len(paths)
	}

	// Simulate the movement of ants until all reach the end
	for {
		activeAnts := false
		var moves []string
		occupiedRooms := make(map[string]bool)
		occupiedTunnels := make(map[string]bool)

		for i := range ants {
			if ants[i].done {
				continue
			}

			path := paths[ants[i].pathIdx]
			currentPos, nextPos := ants[i].position, ants[i].position+1

			// Ensure that the ant is still within the path
			if nextPos < len(path) {
				currentRoom, nextRoom := path[currentPos], path[nextPos]
				tunnelKey, reverseTunnelKey := currentRoom+"-"+nextRoom, nextRoom+"-"+currentRoom

				// Simplified condition for moving to the next room and tunnel
				if !(nextRoom != path[len(path)-1] && occupiedRooms[nextRoom]) &&
					!occupiedTunnels[tunnelKey] && !occupiedTunnels[reverseTunnelKey] {

					// Move the ant to the next position
					ants[i].position = nextPos
					moves = append(moves, fmt.Sprintf("L%d-%s", ants[i].id, nextRoom))

					// Mark room and tunnel as occupied
					if nextRoom != path[len(path)-1] {
						occupiedRooms[nextRoom] = true
					}
					occupiedTunnels[tunnelKey] = true
					occupiedTunnels[reverseTunnelKey] = true

					activeAnts = true
					if nextPos == len(path)-1 {
						ants[i].done = true
					}
				}
			}
		}

		// Print the moves for the current turn
		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}

		// Break the loop if no ants moved
		if !activeAnts {
			break
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <input_file>")
		return
	}

	limits := PathLimits{
		MaxRoomsPerPath:       15,
		MaxTotalPaths:         100,
		MaxPathCombination:    20,
		MaxDirectPaths:        3,
		MaxPathsInCombination: 20,
	}

	// Only print farm information and continue with simulation if valid file and paths exist
	farm, err := loadFarm(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid file format: %v\n", err)
		return
	}
	paths := findAllPaths(farm.AdjacencyList, farm.Start, farm.End, limits)
	if len(paths) == 0 {
		fmt.Println("No paths found from start to end.")
		return
	}

	fmt.Printf("Number of ants: %d\n", farm.AntCount)
	fmt.Printf("Start: %s\n", farm.Start)
	fmt.Printf("End: %s\n", farm.End)
	fmt.Printf("Other rooms: ")
	hasOtherRooms := false
	for room := range farm.AdjacencyList {
		if room != farm.Start && room != farm.End {
			fmt.Printf("%s ", room)
			hasOtherRooms = true
		}
	}
	if !hasOtherRooms {
		fmt.Printf("no other rooms found")
	}
	fmt.Printf("\n\n")

	farm.Paths = paths
	optimalPaths, optimalAntQueue := findOptimalPathCombination(paths, farm.AntCount, limits)

	fmt.Println("Path Combination Selected:")
	for i, path := range optimalPaths {
		fmt.Printf("Path %d: %v\n", i+1, path)
	}
	fmt.Println()

	// Calculate total turns before simulation
	pathLengths := make([]int, len(optimalPaths))
	for i := range optimalPaths {
		pathLengths[i] = len(optimalPaths[i]) - 1
	}
	totalTurns := maxPathTurns(pathLengths, optimalAntQueue)

	simulateAntMovement(optimalPaths, optimalAntQueue, farm.AntCount)
	fmt.Printf("\nNumber of Turns = %d\n", totalTurns)
}
