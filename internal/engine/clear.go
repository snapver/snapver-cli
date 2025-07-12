package engine

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ClearSession deletes snapver branches and diff data.
// If `all` is true, deletes all snapver- branches and all .snapver data.
// If `all` is false, deletes only the current snapver- branch and its diff data.
// defaultBranch: branch to checkout before deleting snapver branches
func ClearSession(all bool, defaultBranch string) {
	if all {
		if err := DeleteAllSnapverBranches(defaultBranch); err != nil {
			fmt.Println("Failed to delete snapver- branches:", err)
		}
		fmt.Println("Removing all .snapver data...")
		if err := removeAllSnapverData(); err != nil {
			fmt.Println("Failed to remove .snapver:", err)
		}
	} else {
		current, err := exec.Command("git", "branch", "--show-current").Output()
		if err == nil {
			branch := strings.TrimSpace(string(current))
			if strings.HasPrefix(branch, "snapver-") {
				if err := DeleteSnapverBranch(branch, defaultBranch); err != nil {
					fmt.Println("Failed to delete branch:", err)
				}
			}
			fmt.Println("Removing diff data for branch:", branch)
			if err := removeBranchDiffData(branch); err != nil {
				fmt.Println("Failed to remove diff dir:", err)
			}
		}
	}
	fmt.Println("Snapver session cleared.")
}

// Delete all snapver- branches in the repository
func DeleteAllSnapverBranches(defaultBranch string) error {
	branches, err := GetBranchList()
	if err != nil {
		return err
	}
	_ = exec.Command("git", "checkout", defaultBranch).Run()
	for _, branch := range branches {
		if strings.HasPrefix(branch, "snapver-") {
			fmt.Println("Deleting branch:", branch)
			_ = exec.Command("git", "branch", "-D", branch).Run()
		}
	}
	return nil
}

// Delete a single snapver- branch by name
func DeleteSnapverBranch(branch, defaultBranch string) error {
	current, _ := exec.Command("git", "branch", "--show-current").Output()
	currentBranch := strings.TrimSpace(string(current))
	if branch == currentBranch {
		_ = exec.Command("git", "checkout", defaultBranch).Run()
	}
	fmt.Println("Deleting branch:", branch)
	return exec.Command("git", "branch", "-D", branch).Run()
}

// Remove all .snapver data
func removeAllSnapverData() error {
	return os.RemoveAll(".snapver")
}

// Remove only the diff data for a specific branch
func removeBranchDiffData(branch string) error {
	diffDir := ".snapver/diffs/" + branch
	return os.RemoveAll(diffDir)
}
