package main

import (
	"bufio"
	"fmt"
	"os"
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

	// debugging: Print file contents
	for _, line := range lines {
		fmt.Println(line)
	}
}
