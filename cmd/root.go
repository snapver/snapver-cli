package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snapver",
	Short: "Automatic file change tracker",
	Long: `
╭─ Snapver ─────────────────────╮
│                               │
│ Automatic file change tracker │
│                               │
╰───────────────────────────────╯
For more info, see: https://github.com/snapver/snapver-cli
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// You can define persistent flags and configuration settings here if needed.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.snapver.yaml)")
}
