package engine

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

// GetBranchList returns a list of all local branches in the repository.
func GetBranchList() ([]string, error) {
	branchesOut, err := exec.Command("git", "branch").Output()
	if err != nil {
		return nil, err
	}
	branches := strings.Split(string(branchesOut), "\n")
	var result []string
	for _, b := range branches {
		branch := strings.TrimSpace(strings.TrimPrefix(b, "* "))
		if branch != "" {
			result = append(result, branch)
		}
	}
	return result, nil
}

// run executes a command and returns error if failed
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// latestCommitHash gets the latest commit hash
func latestCommitHash() (string, error) {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// isIgnored checks whether a file is ignored by git
func isIgnored(path string) bool {
	cmd := exec.Command("git", "check-ignore", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return err == nil && out.Len() > 0
}

// getCurrentBranch gets the current branch name
func getCurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// ensureSnapverInGitignore makes sure .snapver is listed in .gitignore
func ensureSnapverInGitignore(gitignorePath string) {
	const snapverEntry = ".snapver"
	var needWrite bool
	var lines []string

	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		// .gitignore does not exist, create it
		lines = []string{snapverEntry}
		needWrite = true
	} else {
		// .gitignore exists, check if .snapver is present
		content, err := os.ReadFile(gitignorePath)
		if err != nil {
			log.Println("Failed to read .gitignore:", err)
			return
		}
		lines = strings.Split(string(content), "\n")
		found := false
		for _, line := range lines {
			if strings.TrimSpace(line) == snapverEntry {
				found = true
				break
			}
		}
		if !found {
			lines = append(lines, snapverEntry)
			needWrite = true
		}
	}

	if needWrite {
		output := strings.Join(lines, "\n")
		err := os.WriteFile(gitignorePath, []byte(output), 0644)
		if err != nil {
			log.Println("Failed to write .gitignore:", err)
		}
	}
}
