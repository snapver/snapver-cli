package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/snapver/snapver-cli/internal/engine"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start tracking changes",
	Run: func(cmd *cobra.Command, args []string) {
		printControlInfo()

		// check if .git exists
		if !isGitRepo() {
			fmt.Print("No git repo found. Initialize one? (Y/n): ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input == "n" {
				fmt.Println("❌ Aborted.")
				return
			}
			if err := exec.Command("git", "init").Run(); err != nil {
				fmt.Println("❌ git init failed:", err)
				return
			}
			fmt.Println("✓ Initialized git repository.")
		}

		// create .snapver/ directory
		if err := os.MkdirAll(".snapver", 0755); err != nil {
			fmt.Println("❌ Failed to create .snapver directory:", err)
			return
		}
		fmt.Println("✓ Created .snapver/ directory.")

		// create UTC branch
		utc := time.Now().UTC().Format("20060102T150405Z")
		branchName := "snapver-" + utc
		if err := exec.Command("git", "checkout", "-b", branchName).Run(); err != nil {
			fmt.Println("❌ Failed to create branch:", err)
			return
		}
		fmt.Println("✓ Created and switched to branch:", branchName)

		// initial commit
		if err := exec.Command("git", "add", "-A").Run(); err != nil {
			fmt.Println("❌ git add failed:", err)
			return
		}

		if !hasChanges() {
			fmt.Println("No changes to commit → skipping initial commit.")
		} else {
			if err := exec.Command("git", "commit", "-m", "snapver: initial commit").Run(); err != nil {
				fmt.Println("❌ git commit failed:", err)
				return
			}
			fmt.Println("✓ Initial commit created.")
		}

		// Start in background
		d, err := engine.NewEngine(".")
		if err != nil {
			fmt.Println("❌ Failed to start daemon:", err)
			return
		}
		d.StartInBackground()

		// Listen for keypresses
		if err := keyboard.Open(); err != nil {
			fmt.Println("❌ Failed to open keyboard:", err)
			return
		}
		defer keyboard.Close()

		for {
			char, _, err := keyboard.GetSingleKey()
			if err != nil {
				fmt.Println("❌ Keyboard error:", err)
				break
			}
			if char == 'q' {
				fmt.Println("👋 Exiting Snapver")
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func isGitRepo() bool {
	info, err := os.Stat(".git")
	return err == nil && info.IsDir()
}

func hasChanges() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	return err != nil // true = changed
}

func printControlInfo() {
	fmt.Print(`
[ Controls ]
q : Exit Snapver

`)
}
