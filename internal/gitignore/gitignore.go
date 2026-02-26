package gitignore

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Entry struct {
	Name    string
	Content string
}

func Add(path, templateName, content string, override bool) error {
	return AddMultiple(path, []Entry{{Name: templateName, Content: content}}, override)
}

func AddMultiple(path string, entries []Entry, override bool) error {
	var sections []string
	for _, e := range entries {
		sections = append(sections, wrapWithMarkers(e.Name, e.Content))
	}
	newContent := strings.Join(sections, "\n\n")

	if override {
		return os.WriteFile(path, []byte(newContent+"\n"), 0644)
	}

	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(path, []byte(newContent+"\n"), 0644)
		}
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	result := string(existing)
	for _, e := range entries {
		wrapped := wrapWithMarkers(e.Name, e.Content)
		pattern := regexp.MustCompile(
			fmt.Sprintf(`(?s)# >>> ignore-cli: %s\n.*?# <<< ignore-cli: %s\n?`,
				regexp.QuoteMeta(e.Name), regexp.QuoteMeta(e.Name)))

		if pattern.MatchString(result) {
			result = pattern.ReplaceAllString(result, wrapped+"\n")
		} else {
			result = strings.TrimRight(result, "\n") + "\n\n" + wrapped + "\n"
		}
	}

	return os.WriteFile(path, []byte(result), 0644)
}

func wrapWithMarkers(name, content string) string {
	start := fmt.Sprintf("# >>> ignore-cli: %s", name)
	end := fmt.Sprintf("# <<< ignore-cli: %s", name)
	return fmt.Sprintf("%s\n%s\n%s", start, strings.TrimSpace(content), end)
}
