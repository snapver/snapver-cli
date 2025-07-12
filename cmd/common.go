package cmd

import (
	"bufio"
	"fmt"
	"strings"
)

// selectDefaultBranch prompts the user to select a default branch according to the rules.
func selectDefaultBranch(reader *bufio.Reader, branches []string) (string, bool) {
	mainExists, masterExists := false, false
	for _, b := range branches {
		if b == "main" {
			mainExists = true
		}
		if b == "master" {
			masterExists = true
		}
	}
	if mainExists && masterExists {
		fmt.Println("Both 'main' and 'master' branches exist.")
		fmt.Print("Which branch do you want to use as default? (main/master): ")
		resp, _ := reader.ReadString('\n')
		resp = strings.TrimSpace(strings.ToLower(resp))
		if resp == "main" || resp == "master" {
			return resp, true
		}
		fmt.Println("Invalid input. Aborted.")
		return "", false
	}
	if mainExists {
		return "main", true
	}
	if masterExists {
		return "master", true
	}
	fmt.Println("No 'main' or 'master' branch found.")
	fmt.Println("Available branches:")
	for i, b := range branches {
		fmt.Printf("  [%d] %s\n", i+1, b)
	}
	fmt.Print("Select default branch by number: ")
	resp, _ := reader.ReadString('\n')
	resp = strings.TrimSpace(resp)
	idx := -1
	fmt.Sscanf(resp, "%d", &idx)
	if idx < 1 || idx > len(branches) {
		fmt.Println("Invalid selection. Aborted.")
		return "", false
	}
	return branches[idx-1], true
}
