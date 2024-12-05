package farm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Farm struct {
	AntCount      int
	Start         string
	End           string
	AdjacencyList map[string][]string
	Paths         [][]string
}

func Load(filename string) (*Farm, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("empty input file")
	}

	antCount, err := strconv.Atoi(lines[0])
	if err != nil || antCount <= 0 {
		return nil, fmt.Errorf("invalid number of ants: must be a number greater than 0")
	}

	var start, end string
	var links []string

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		switch line {
		case "##start":
			if i+1 >= len(lines) {
				return nil, fmt.Errorf("invalid file format: missing room definition after start command")
			}
			start = strings.Fields(lines[i+1])[0]
			i++
		case "##end":
			if i+1 >= len(lines) {
				return nil, fmt.Errorf("invalid file format: missing room definition after end command")
			}
			end = strings.Fields(lines[i+1])[0]
			i++
		default:
			if strings.Contains(line, "-") {
				links = append(links, line)
			}
		}
	}

	if start == "" {
		return nil, fmt.Errorf("no start room defined")
	}
	if end == "" {
		return nil, fmt.Errorf("no end room defined")
	}
	if len(links) == 0 {
		return nil, fmt.Errorf("no room connections found")
	}

	adjacencyList, err := buildAdjacencyList(links)
	if err != nil {
		return nil, err
	}
	return &Farm{AntCount: antCount, Start: start, End: end, AdjacencyList: adjacencyList}, nil
}

func buildAdjacencyList(edges []string) (map[string][]string, error) {
	adjList := make(map[string][]string)
	for _, edge := range edges {
		parts := strings.Split(edge, "-")
		if len(parts) != 2 || parts[0] == parts[1] {
			return nil, fmt.Errorf("invalid link between rooms: %s", edge)
		}
		adjList[parts[0]] = append(adjList[parts[0]], parts[1])
		adjList[parts[1]] = append(adjList[parts[1]], parts[0])
	}
	return adjList, nil
}

func (f *Farm) PrintInfo() {
	fmt.Printf("Number of ants: %d\n", f.AntCount)
	fmt.Printf("Start: %s\n", f.Start)
	fmt.Printf("End: %s\n", f.End)
	fmt.Printf("Other rooms: ")
	
	hasOtherRooms := false
	for room := range f.AdjacencyList {
		if room != f.Start && room != f.End {
			fmt.Printf("%s ", room)
			hasOtherRooms = true
		}
	}
	if !hasOtherRooms {
		fmt.Printf("no other rooms found")
	}
	fmt.Printf("\n\n")
}
