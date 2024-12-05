package main

import (
	"fmt"
	"lem-in/pkg/farm"
	"lem-in/pkg/pathfinder"
	"lem-in/pkg/simulator"
	"os"
	"time"
)

func main() {
	startTime := time.Now()
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <input_file>")
		return
	}

	limits := pathfinder.PathLimits{
		MaxRoomsPerPath:       15,
		MaxTotalPaths:         100,
		MaxPathCombination:    20,
		MaxDirectPaths:        3,
		MaxPathsInCombination: 20,
	}

	antFarm, err := farm.Load(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid file format: %v\n", err)
		return
	}

	paths := pathfinder.FindAllPaths(antFarm.AdjacencyList, antFarm.Start, antFarm.End, limits)
	if len(paths) == 0 {
		fmt.Println("No paths found from start to end.")
		return
	}

	antFarm.PrintInfo()
	antFarm.Paths = paths

	optimalPaths, optimalAntQueue := pathfinder.FindOptimalPathCombination(paths, antFarm.AntCount, limits)
	pathfinder.PrintSelectedPaths(optimalPaths)

	turns := simulator.Run(optimalPaths, optimalAntQueue, antFarm.AntCount)
	fmt.Printf("\nNumber of Turns = %d\n", turns)

	executionTime := time.Since(startTime)
	fmt.Printf("Execution time: %v\n", executionTime)

}
