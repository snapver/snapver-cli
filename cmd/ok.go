package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/snapver/snapver-cli/internal/engine"
	"github.com/spf13/cobra"
)

var okCmd = &cobra.Command{
	Use:   "ok",
	Short: "Merge changes into default branch",
	Run: func(cmd *cobra.Command, args []string) {
		current, err := exec.Command("git", "branch", "--show-current").Output()
		if err != nil {
			fmt.Println("❌ Failed to get current branch:", err)
			return
		}
		currentBranch := strings.TrimSpace(string(current))

		if !strings.HasPrefix(currentBranch, "snapver-") {
			fmt.Println("⚠️ Not in a snapver branch:", currentBranch)
			return
		}

		branches, err := engine.GetBranchList()
		if err != nil || len(branches) == 0 {
			fmt.Println("❌ No branches found or failed to get branch list.")
			return
		}
		reader := bufio.NewReader(os.Stdin)
		defaultBranch, ok := selectDefaultBranch(reader, branches)
		if !ok {
			fmt.Println("❌ No main or master branch found. Please create one.")
			return
		}

		fmt.Printf("You are about to merge '%s' into '%s'. Proceed? [y/N]: ", currentBranch, defaultBranch)
		resp, _ := reader.ReadString('\n')
		resp = strings.ToLower(strings.TrimSpace(resp))
		if resp != "y" && resp != "yes" {
			fmt.Println("Aborted.")
			return
		}

		fmt.Println("Merging", currentBranch, "→", defaultBranch)

		usedStash := false
		if err := exec.Command("git", "switch", defaultBranch).Run(); err != nil {
			fmt.Println("❌ Failed to switch to", defaultBranch, ":", err)
			fmt.Println("You may have uncommitted changes. Commit or stash them, or let Snapver auto-stash.")
			fmt.Print("Auto-stash and continue? [y/N]: ")
			resp, _ := reader.ReadString('\n')
			resp = strings.ToLower(strings.TrimSpace(resp))
			if resp != "y" && resp != "yes" {
				fmt.Println("Aborted.")
				return
			}
			if err := exec.Command("git", "stash", "--include-untracked").Run(); err != nil {
				fmt.Println("❌ Failed to stash changes:", err)
				return
			}
			fmt.Println("Changes stashed. Retrying switch...")
			if err := exec.Command("git", "switch", defaultBranch).Run(); err != nil {
				fmt.Println("❌ Still failed to switch to", defaultBranch, ":", err)
				return
			}
			fmt.Println("You can restore your changes with 'git stash pop' after merge.")
			usedStash = true
		}
		if err := exec.Command("git", "merge", "--no-ff", currentBranch).Run(); err != nil {
			fmt.Println("❌ Merge failed:", err)
			return
		}

		fmt.Printf("Delete branch '%s'? [y/N]: ", currentBranch)
		resp, _ = reader.ReadString('\n')
		resp = strings.ToLower(strings.TrimSpace(resp))
		if resp == "y" || resp == "yes" {
			if err := exec.Command("git", "branch", "-d", currentBranch).Run(); err != nil {
				fmt.Println("⚠️ Could not delete branch with -d, trying -D ...")
				if err := exec.Command("git", "branch", "-D", currentBranch).Run(); err != nil {
					fmt.Println("❌ Failed to delete branch:", err)
					return
				}
			}
			fmt.Println("Deleted branch:", currentBranch)
		} else {
			fmt.Println("Branch not deleted. You may want to delete it manually later.")
		}

		if usedStash {
			fmt.Println("Restoring stashed changes with 'git stash pop' ...")
			if err := exec.Command("git", "stash", "pop").Run(); err != nil {
				fmt.Println("⚠️ Failed to pop stash. You may need to resolve conflicts manually.")
			} else {
				fmt.Println("Stashed changes restored.")
			}
		}
		fmt.Println("✅ Merged into", defaultBranch)
	},
}

func init() {
	rootCmd.AddCommand(okCmd)
}
