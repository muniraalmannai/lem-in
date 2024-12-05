package pathfinder

import (
	"fmt"
	"sort"
)

type PathLimits struct {
	MaxRoomsPerPath       int
	MaxTotalPaths         int
	MaxPathCombination    int
	MaxDirectPaths        int
	MaxPathsInCombination int
}

func NewPathLimits() PathLimits {
	return PathLimits{
		MaxRoomsPerPath:       15,
		MaxTotalPaths:         100,
		MaxPathCombination:    20,
		MaxDirectPaths:        3,
		MaxPathsInCombination: 20,
	}
}

func FindAllPaths(adjacencyList map[string][]string, start, end string, limits PathLimits) [][]string {
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

	return paths
}

func contains(path []string, room string) bool {
	for _, v := range path {
		if v == room {
			return true
		}
	}
	return false
}

func FindOptimalPathCombination(paths [][]string, antCount int, limits PathLimits) ([][]string, []int) {
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	if len(paths) > limits.MaxPathCombination {
		paths = paths[:limits.MaxPathCombination]
	}

	if containsDirectPath(paths) {
		return handleDirectPaths(paths, antCount, limits)
	}

	return handleComplexPaths(paths, antCount, limits)
}

func containsDirectPath(paths [][]string) bool {
	for _, path := range paths {
		if len(path) == 2 {
			return true
		}
	}
	return false
}

func handleDirectPaths(paths [][]string, antCount int, limits PathLimits) ([][]string, []int) {
	var optimalPaths [][]string
	for _, path := range paths {
		if isValidCombination(append(optimalPaths, path)) {
			optimalPaths = append(optimalPaths, path)
		}
		if len(optimalPaths) >= limits.MaxDirectPaths {
			break
		}
	}

	pathLengths := make([]int, len(optimalPaths))
	for i := range optimalPaths {
		pathLengths[i] = len(optimalPaths[i]) - 1
	}
	
	return optimalPaths, distributeAnts(optimalPaths, antCount)
}

func handleComplexPaths(paths [][]string, antCount int, limits PathLimits) ([][]string, []int) {
	var optimalPaths [][]string
	var optimalAntQueue []int
	minTurns := int(^uint(0) >> 1)

	for i := 1; i <= len(paths); i++ {
		combinations := generateCombinations(paths, i, limits)
		for _, combo := range combinations {
			if isValidCombination(combo) {
				antQueue, turns := calculateDistribution(combo, antCount)
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

func isValidCombination(paths [][]string) bool {
	roomSet := make(map[string]bool)
	for _, path := range paths {
		for _, room := range path[1 : len(path)-1] {
			if roomSet[room] {
				return false
			}
			roomSet[room] = true
		}
	}
	return true
}

func generateCombinations(paths [][]string, length int, limits PathLimits) [][][]string {
	if length > limits.MaxPathsInCombination {
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

func calculateDistribution(paths [][]string, antCount int) ([]int, int) {
	antQueue := make([]int, len(paths))
	pathLengths := make([]int, len(paths))
	for i := range paths {
		pathLengths[i] = len(paths[i]) - 1
	}

	remaining := antCount
	for remaining > 0 {
		shortestIdx := findShortestPath(pathLengths, antQueue)
		antQueue[shortestIdx]++
		remaining--
	}

	return antQueue, maxPathTurns(pathLengths, antQueue)
}

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

func PrintSelectedPaths(paths [][]string) {
	fmt.Println("Path Combination Selected:")
	for i, path := range paths {
		fmt.Printf("Path %d: %v\n", i+1, path)
	}
	fmt.Println()
}

func distributeAnts(paths [][]string, antCount int) []int {
    antQueue := make([]int, len(paths))
    pathLengths := make([]int, len(paths))
    for i := range paths {
        pathLengths[i] = len(paths[i]) - 1
    }

    remaining := antCount
    for remaining > 0 {
        shortestIdx := findShortestPath(pathLengths, antQueue)
        antQueue[shortestIdx]++
        remaining--
    }

    return antQueue
}
