package simulator

import (
	"fmt"
	"strings"
)

type ant struct {
	id       int
	pathIdx  int
	position int
	done     bool
}

func Run(paths [][]string, antQueue []int, antCount int) int {
	ants := make([]ant, antCount)
	antId, pathIdx := 1, 0

	// Initialize ants and assign them to paths
	for antId <= antCount {
		if antQueue[pathIdx] > 0 {
			ants[antId-1] = ant{id: antId, pathIdx: pathIdx, position: 0}
			antQueue[pathIdx]--
			antId++
		}
		pathIdx = (pathIdx + 1) % len(paths)
	}

	turns := 0
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

			if nextPos < len(path) {
				currentRoom, nextRoom := path[currentPos], path[nextPos]
				tunnelKey := currentRoom + "-" + nextRoom
				reverseTunnelKey := nextRoom + "-" + currentRoom

				if !(nextRoom != path[len(path)-1] && occupiedRooms[nextRoom]) &&
					!occupiedTunnels[tunnelKey] && !occupiedTunnels[reverseTunnelKey] {

					ants[i].position = nextPos
					moves = append(moves, fmt.Sprintf("L%d-%s", ants[i].id, nextRoom))

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

		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
			turns++
		}

		if !activeAnts {
			break
		}
	}

	return turns
}
