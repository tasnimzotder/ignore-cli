package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "ignore",
	Short: "A CLI tool to manage .gitignore files",
	Long:  `ignore is a CLI tool to manage .gitignore files. It allows you to list, search, and add .gitignore templates to your project.`,
}

func init() {
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(searchCmd)
	RootCmd.AddCommand(addCmd)
}
