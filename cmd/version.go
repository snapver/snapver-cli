package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of snapver
	Version = "v0.0.0" // fallback
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display version information for snapver.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Snapver %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
