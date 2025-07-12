package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop tracking",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("pgrep", "-f", "snapver start").Output()
		if err != nil {
			fmt.Println("No running Snapver found.")
			return
		}
		pids := strings.FieldsSeq(string(out))
		for pid := range pids {
			exec.Command("kill", pid).Run()
			fmt.Println("Killed Snapver:", pid)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
