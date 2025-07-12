package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/snapver/snapver-cli/internal/engine"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete snapver branches and diff data.",
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		reader := bufio.NewReader(os.Stdin)
		branches, err := engine.GetBranchList()
		if err != nil || len(branches) == 0 {
			fmt.Println("No branches found or failed to get branch list.")
			return
		}

		defaultBranch, ok := selectDefaultBranch(reader, branches)
		if !ok {
			return
		}

		if all {
			fmt.Print("This will delete ALL snapver- branches and all .snapver data. Are you sure? [y/N]: ")
		} else {
			fmt.Print("This will delete the current snapver- branch and its diff data. Are you sure? [y/N]: ")
		}
		resp, _ := reader.ReadString('\n')
		resp = strings.TrimSpace(strings.ToLower(resp))
		if resp != "y" && resp != "yes" {
			fmt.Println("Aborted.")
			return
		}

		engine.ClearSession(all, defaultBranch)
	},
}

func init() {
	clearCmd.Flags().Bool("all", false, "Remove all .snapver data (not just current branch)")
	rootCmd.AddCommand(clearCmd)
}
