package cmd

import (
	"sort"

	"github.com/spf13/cobra"
	"github.com/tasnimzotder/ignore-cli/internal/search"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available .gitignore templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		templates, err := search.AllTemplates()
		if err != nil {
			return err
		}

		sort.Strings(templates)

		for _, template := range templates {
			cmd.Println(template)
		}

		return nil
	},
}
