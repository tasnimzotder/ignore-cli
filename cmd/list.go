package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Browse available .gitignore templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitignorePath, _ := filepath.Abs(".gitignore")
		return launchTUI(gitignorePath)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
