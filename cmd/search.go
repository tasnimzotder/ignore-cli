package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tasnimzotder/ignore-cli/internal/template"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for .gitignore templates",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		results, err := template.Search(args[0])
		if err != nil {
			return err
		}

		for _, result := range results {
			cmd.Println(result.Name)
		}

		return nil
	},
}
