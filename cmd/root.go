package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ignore",
	Short: "A TUI tool to manage .gitignore files",
	Long:  `ignore is a CLI tool to manage .gitignore files using GitHub's gitignore templates. Features an interactive template picker with fuzzy search.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
