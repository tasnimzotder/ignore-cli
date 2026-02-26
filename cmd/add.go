package cmd

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/tasnimzotder/ignore-cli/internal/cache"
	"github.com/tasnimzotder/ignore-cli/internal/gitignore"
	"github.com/tasnimzotder/ignore-cli/internal/tui"
)

var overrideFlag bool

var addCmd = &cobra.Command{
	Use:   "add [templates...]",
	Short: "Add .gitignore templates (interactive picker if no args)",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitignorePath, _ := filepath.Abs(".gitignore")

		if len(args) > 0 {
			return addDirect(gitignorePath, args)
		}
		return addInteractive(gitignorePath)
	},
}

func init() {
	addCmd.Flags().BoolVarP(&overrideFlag, "override", "o", false, "Replace existing .gitignore instead of appending")
	rootCmd.AddCommand(addCmd)
}

func addDirect(path string, templateNames []string) error {
	c := cache.GetInstance()

	var entries []gitignore.Entry
	for _, name := range templateNames {
		tmpl, err := c.FindTemplate(name)
		if err != nil {
			return fmt.Errorf("template %q: %w", name, err)
		}
		content, err := tmpl.Content()
		if err != nil {
			return fmt.Errorf("template %q: %w", name, err)
		}
		entries = append(entries, gitignore.Entry{Name: tmpl.Name, Content: content})
	}

	if err := gitignore.AddMultiple(path, entries, overrideFlag); err != nil {
		return err
	}

	for _, e := range entries {
		fmt.Printf("Added %s template\n", e.Name)
	}
	return nil
}

func addInteractive(path string) error {
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

	picker := tui.NewPicker(fetchFn)
	p := tea.NewProgram(picker, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	result := finalModel.(tui.PickerModel).Result()
	if result.Canceled || len(result.Selected) == 0 {
		fmt.Println("No templates selected.")
		return nil
	}

	var entries []gitignore.Entry
	for _, name := range result.Selected {
		tmpl, err := c.FindTemplate(name)
		if err != nil {
			return fmt.Errorf("template %q: %w", name, err)
		}
		content, err := tmpl.Content()
		if err != nil {
			return fmt.Errorf("template %q: %w", name, err)
		}
		entries = append(entries, gitignore.Entry{Name: tmpl.Name, Content: content})
	}

	if err := gitignore.AddMultiple(path, entries, overrideFlag); err != nil {
		return err
	}

	for _, e := range entries {
		fmt.Printf("Added %s template\n", e.Name)
	}
	return nil
}
