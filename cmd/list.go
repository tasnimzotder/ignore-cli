package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tasnimzotder/ignore-cli/internal/template"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available .gitignore templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		templates, err := template.List()
		if err != nil {
			return err
		}

		for _, t := range templates {
			cmd.Println(t.Name)
		}

		return nil
	},
}
