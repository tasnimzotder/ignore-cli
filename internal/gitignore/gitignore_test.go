package gitignore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd_CreatesNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	err := Add(path, "Go", "*.exe\n*.dll", false)
	require.NoError(t, err)

	content, _ := os.ReadFile(path)
	got := string(content)

	assert.Contains(t, got, "# >>> ignore-cli: Go")
	assert.Contains(t, got, "# <<< ignore-cli: Go")
	assert.Contains(t, got, "*.exe")
}

func TestAdd_AppendsToExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	os.WriteFile(path, []byte("# existing\nnode_modules/\n"), 0644)

	err := Add(path, "Go", "*.exe\n*.dll", false)
	require.NoError(t, err)

	content, _ := os.ReadFile(path)
	got := string(content)

	assert.Contains(t, got, "node_modules/")
	assert.Contains(t, got, "# >>> ignore-cli: Go")
}

func TestAdd_ReplacesExistingSection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	existing := "# >>> ignore-cli: Go\n*.exe\n# <<< ignore-cli: Go\n"
	os.WriteFile(path, []byte(existing), 0644)

	err := Add(path, "Go", "*.exe\n*.dll\n*.so", false)
	require.NoError(t, err)

	content, _ := os.ReadFile(path)
	got := string(content)

	assert.Contains(t, got, "*.so")
	assert.Equal(t, 1, strings.Count(got, "# >>> ignore-cli: Go"))
}

func TestAdd_Override(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	os.WriteFile(path, []byte("# old stuff\nnode_modules/\n"), 0644)

	err := Add(path, "Go", "*.exe", true)
	require.NoError(t, err)

	content, _ := os.ReadFile(path)
	got := string(content)

	assert.NotContains(t, got, "node_modules")
	assert.Contains(t, got, "*.exe")
}

func TestAddMultiple(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	entries := []Entry{
		{Name: "Go", Content: "*.exe\n*.dll"},
		{Name: "Python", Content: "*.pyc\n__pycache__/"},
	}

	err := AddMultiple(path, entries, false)
	require.NoError(t, err)

	content, _ := os.ReadFile(path)
	got := string(content)

	assert.Contains(t, got, "# >>> ignore-cli: Go")
	assert.Contains(t, got, "# >>> ignore-cli: Python")
	assert.Contains(t, got, "*.exe")
	assert.Contains(t, got, "*.pyc")
}

func TestAddMultiple_Override(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	os.WriteFile(path, []byte("old content\n"), 0644)

	entries := []Entry{
		{Name: "Go", Content: "*.exe"},
		{Name: "Python", Content: "*.pyc"},
	}

	err := AddMultiple(path, entries, true)
	require.NoError(t, err)

	content, _ := os.ReadFile(path)
	got := string(content)

	assert.NotContains(t, got, "old content")
	assert.Contains(t, got, "# >>> ignore-cli: Go")
	assert.Contains(t, got, "# >>> ignore-cli: Python")
}
