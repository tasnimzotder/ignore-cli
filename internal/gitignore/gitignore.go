package gitignore

import (
	"fmt"
	"os"

	"github.com/tasnimzotder/ignore-cli/internal/search"
	"github.com/tasnimzotder/ignore-cli/internal/template"
)

func Add(projectType string, override bool) error {
	// if cache is empty, fetch the list of templates
	_, err := search.AllTemplates()
	if err != nil {
		return err
	}

	t, err := template.Get(projectType)
	if err != nil {
		return err
	}

	content, err := t.Content()
	if err != nil {
		return err
	}

	return writeGitignore(content, override)
}

func writeGitignore(content string, override bool) error {
	if !override {
		existing, err := os.ReadFile(".gitignore")
		if err == nil {
			content = string(existing) + "\n" + content
		}
	}

	err := os.WriteFile(".gitignore", []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write .gitignore file: %w", err)
	}

	fmt.Println("Successfully added/updated .gitignore file.")
	return nil
}
