package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	Name  string
	Jeran []*Room
}
type Farm struct {
	Rooms  map[string]*Room
	AntNum int
	Start  Room
	End    Room
}

var (
	farm        Farm
	mwjoodStart bool
	mwjoodEnd   bool
	allPaths    [][]string
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("go run . example.txt")
		return
	}
	ParseFile(os.Args[1])
	Duffs(farm.Start.Name, farm.End.Name, []string{}, &allPaths)
	if len(allPaths) == 0 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(0)
	}
	var result1 [][][]string
	generateSubsets(allPaths, 0, [][]string{}, &result1)
	bestPaths := ChoiseCollectionPaths(result1)
	MoveAnts(bestPaths, farm.AntNum)
}
func ParseFile(file string) {
	open, or := os.ReadFile(file)
	Err(or)
	split := strings.Split(string(open), "\n")
	farm.AntNum, or = strconv.Atoi(strings.TrimSpace(split[0]))
	if farm.AntNum == 0 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(0)
	}
	Err(or)
	farm.Rooms = make(map[string]*Room)
	for _, line := range split[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		f := strings.Fields(line)
		if line == "##start" {
			mwjoodStart = true
			continue
		}
		if line == "##end" {
			mwjoodEnd = true
			continue
		}
		if mwjoodStart {
			farm.Start.Name = f[0]
			mwjoodStart = false
			continue
		}
		if mwjoodEnd {
			farm.End.Name = f[0]
			mwjoodEnd = false
			continue
		}
		if strings.Contains(line, "-") {
			split := strings.Split(line, "-")
			from := split[0]
			to := split[1]
			mapmaker(from, to)
		}
	}
}
func Err(Error error) {
	if Error != nil {
		fmt.Println(Error)
		os.Exit(0)
	}
}
func mapmaker(from, to string) {
	fromNode, fromExists := farm.Rooms[from]
	if !fromExists {
		fromNode = &Room{Name: from}
		farm.Rooms[from] = fromNode
	}
	toNode, toExists := farm.Rooms[to]
	if !toExists {
		toNode = &Room{Name: to}
		farm.Rooms[to] = toNode
	}
	fromNode.Jeran = append(fromNode.Jeran, toNode)
	toNode.Jeran = append(toNode.Jeran, fromNode)
}
func Duffs(start, end string, path []string, allPaths *[][]string) {
	path = append(path, start)
	if start == end {
		pathCopy := make([]string, len(path))
		copy(pathCopy, path)
		*allPaths = append(*allPaths, pathCopy)
		return
	}
	StartRoom := farm.Rooms[start]
	for _, neighbor := range StartRoom.Jeran {
		visited := false
		for _, node := range path {
			if neighbor.Name == node {
				visited = true
				break
			}
		}
		if !visited {
			Duffs(neighbor.Name, end, path, allPaths)
		}
	}
}

// Function to generate all non-empty subsets of allPaths
func generateSubsets(paths [][]string, index int, currentSubset [][]string, result *[][][]string) {
	if index == len(paths) {
		// If the current subset is non-empty, add it to the result
		if len(currentSubset) > 0 && len(currentSubset) <= farm.AntNum {
			// Make a deep copy of currentSubset and add it to the result
			subsetCopy := make([][]string, len(currentSubset))
			for i := range currentSubset {
				subsetCopy[i] = make([]string, len(currentSubset[i]))
				copy(subsetCopy[i], currentSubset[i])
			}
			*result = append(*result, subsetCopy)
		}
		return
	}

	// Option 1: Exclude the current path and move to the next
	generateSubsets(paths, index+1, currentSubset, result)

	// Option 2: Include the current path and move to the next
	newSubset := append(currentSubset, paths[index])
	if !Conflicts(newSubset) {
		generateSubsets(paths, index+1, newSubset, result)
	}
}

// Check if there are conflicts between any two paths in the subset
func Conflicts(result [][]string) bool {
	for index1, path1 := range result {
		for index2 := index1 + 1; index2 < len(result); index2++ {
			path2 := result[index2]
			if Resolve(path1, path2) {
				return true
			}
		}
	}
	return false
}

// Check if two paths conflict (i.e., share any rooms excluding start/end)
func Resolve(path1, path2 []string) bool {
	for _, room1 := range path1[1 : len(path1)-1] {
		for _, room2 := range path2[1 : len(path2)-1] {
			if room1 == room2 {
				return true
			}
		}
	}
	return false
}

func ChoiseCollectionPaths(grpPaths [][][]string) [][]string {
	bestPaths := grpPaths[0]
	num := (farm.AntNum / len(bestPaths)) * countAvgNumRooms(bestPaths)
	for _, grp := range grpPaths[1:] {
		if len(grp) <= farm.AntNum {
			equationRes := (farm.AntNum / len(grp)) * countAvgNumRooms(grp)
			if equationRes < num {
				bestPaths = grp
				num = equationRes
			}
		}
	}
	return bestPaths
}
func countAvgNumRooms(grp [][]string) int {
	numRooms := 0
	for _, path := range grp {
		numRooms = numRooms + len(path)
	}
	return (numRooms / len(grp))
}

func MoveAnts(paths [][]string, numAnts int) {
	// Array to track how many ants are in each path
	antsInPaths := make([]int, len(paths))

	// Ant positions and room occupation
	antPositions := make([]int, numAnts)
	roomOccupied := make(map[string]int)

	// Track ant movements
	antAssignments := make([]int, numAnts)

	// To track if a 2-room path is occupied in the current move
	shortPathUsed := make([]bool, len(paths))

	// Assign ants to paths based on the rooms + ants logic
	for i := 0; i < numAnts; i++ {
		bestPath := 0
		bestValue := len(paths[0]) + antsInPaths[0]
		for j := 1; j < len(paths); j++ {
			currentValue := len(paths[j]) + antsInPaths[j]
			if currentValue < bestValue {
				bestPath = j
				bestValue = currentValue
			}
		}
		// Assign ant to the best path
		antAssignments[i] = bestPath
		antsInPaths[bestPath]++
	}

	// Now move ants through the paths
	for move := 0; ; move++ {
		var movements []string
		allAntsMoved := true

		// Reset shortPathUsed for each move
		for i := range shortPathUsed {
			shortPathUsed[i] = false
		}

		for i := 0; i < numAnts; i++ {
			path := paths[antAssignments[i]]
			pos := antPositions[i]

			if pos < len(path)-1 {
				nextPos := pos + 1
				nextRoom := path[nextPos]

				// Ensure paths with length 2 (start and end) are used only once in a move
				if len(path) == 2 && shortPathUsed[antAssignments[i]] {
					continue // Skip this path if already used in this move
				}

				if nextRoom == farm.End.Name || roomOccupied[nextRoom] == 0 {
					if pos > 0 {
						roomOccupied[path[pos]]--
					}
					antPositions[i] = nextPos
					roomOccupied[nextRoom]++
					movements = append(movements, fmt.Sprintf("L%d-%s", i+1, nextRoom))

					// Mark the short path as used in this move
					if len(path) == 2 {
						shortPathUsed[antAssignments[i]] = true
					}

					allAntsMoved = false
				}
			}
		}

		if allAntsMoved {
			break
		}

		if len(movements) > 0 {
			fmt.Println(strings.Join(movements, " "))
		}
	}
}
