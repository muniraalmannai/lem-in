# 🐜 Ant Farm (lem-in)

![image](images/image.gif)


This project, **lem-in**, is a digital simulation of an ant farm where ants navigate through rooms and tunnels to reach an exit. The objective is to find the quickest way to move a set number of ants from a designated start room (`##start`) to an end room (`##end`) while avoiding congestion and optimizing their paths.

## Project Overview

The program will:
1. Read input from a file describing the ants, rooms, and tunnels of the colony.
2. Parse and validate the input to ensure correct format and logical structure.
3. Calculate the most efficient path(s) for the ants to reach the end room.
4. Display the movement of each ant across the colony on each turn.

### How It Works

The goal of this project is to find the quickest way to get n ants across the colony. There are some constraints:
- The ants need to take the shortest path (which isn't always the simplest).
- They must avoid traffic jams and prevent overcrowding in any room.

To achieve this, the **Edmonds-Karp** algorithm is implemented to find the most optimal path combinations, ensuring that the ants follow the most efficient routes. The simulation then moves ants through the selected paths while respecting the single occupancy per room and per tunnel per turn constraints, ensuring an efficient traversal across the colony.

### Input File Format

The input file should contain:
- The number of ants.
- Definitions of rooms in the format `name coord_x coord_y` (e.g., `Room 1 2`).
- Definitions of tunnels in the format `name1-name2` (e.g., `1-2`).

Commands like `##start` and `##end` designate the start and end rooms. Additional comments can be added, starting with `#`.

### Example Input

```text
##start
1 23 3
2 16 7
# This is a comment
3 16 3
##end
0 9 5
0-4
0-6
1-3
4-3
```

### Output Format

The program outputs:
1. The number of ants.
2. The room definitions.
3. The tunnel definitions.
4. Movements of ants in each turn, formatted as `Lx-y`, where `x` is the ant number and `y` is the room name (e.g., `L1-2` means ant 1 moved to room 2).

### Sample Output

```text
L1-2 L2-3
L1-4 L2-5
```

### Features

- **Shortest Path Calculation:** Ants find the quickest path to `##end` based on optimal route planning, avoiding congestion.
- **Error Handling:** The program validates input and handles errors like missing `##start` or `##end` rooms, duplicate rooms, unknown commands, and formatting issues.
- **Single Occupancy Constraint:** Each room (except `##start` and `##end`) holds only one ant at a time.
- **Turn-Based Movement:** Each ant moves once per turn, using tunnels that connect rooms.

### Constraints

- Rooms do not start with `L` or `#`, and room names do not contain spaces.
- Each tunnel links exactly two rooms.
- No room has more than one tunnel to any other room.
- Ants take the shortest path(s), avoiding traffic jams and backtracking.

### Requirements

- **Language:** Go
- **Packages Allowed:** Only standard Go packages

### Usage

To run the program, provide an input file as an argument:

```bash
go run main.go <path_to_input_file>
```

### Error Messages

If any invalid format or logical inconsistency is detected, the program outputs:

```text
ERROR: invalid data format
```

or, for specific errors:

```text
ERROR: invalid data format, [description of the issue]
```

### Development and Testing

The code follows Go best practices, with unit tests provided for key functionalities. Test files are included to validate both expected behavior and error handling.

## Example Run

With an input file as shown above, the output displays each turn's movements until all ants reach `##end`, demonstrating how efficiently they navigate the colony's layout.

## Authors

- Munira Almannai
- Zahra Mahdi
- Abdulroda Salman

## License

This project is for educational purposes as part of a coursework assignment.
