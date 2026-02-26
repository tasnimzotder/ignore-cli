package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/tasnimzotder/ignore-cli/internal/cache"
	"github.com/tasnimzotder/ignore-cli/internal/tui"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Browse available .gitignore templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := cache.GetInstance()
		fetchFn := func() ([]string, error) {
			templates, err := c.ListTemplates()
			if err != nil {
				return nil, err
			}
			names := make([]string, len(templates))
			for i, t := range templates {
				names[i] = t.Name
			}
			return names, nil
		}

		browser := tui.NewBrowser(fetchFn)
		p := tea.NewProgram(browser, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
