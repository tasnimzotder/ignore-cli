package cache

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetSingleton() {
	instance = nil
	once = sync.Once{}
}

func TestGetInstance_ReturnsSameInstance(t *testing.T) {
	resetSingleton()
	t.Setenv("HOME", t.TempDir())

	a := GetInstance()
	b := GetInstance()
	assert.Same(t, a, b, "GetInstance should return the same pointer")
}

func TestListTemplates_FetchesFromAPI(t *testing.T) {
	resetSingleton()
	t.Setenv("HOME", t.TempDir())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{"Go", "Python", "Rust"})
	}))
	defer server.Close()

	origURL := gitignoreListURL
	gitignoreListURL = server.URL
	defer func() { gitignoreListURL = origURL }()

	c := GetInstance()
	templates, err := c.ListTemplates()
	require.NoError(t, err)
	assert.Len(t, templates, 3)
	assert.Equal(t, "Go", templates[0].Name)
}

func TestSearchTemplates_CaseInsensitive(t *testing.T) {
	resetSingleton()
	t.Setenv("HOME", t.TempDir())

	c := GetInstance()
	c.Templates = []Template{
		{Name: "Go"},
		{Name: "Python"},
		{Name: "Node"},
		{Name: "Go.AllowList"},
	}

	results := c.SearchTemplates("go")
	assert.Len(t, results, 2)
	assert.Equal(t, "Go", results[0].Name)
	assert.Equal(t, "Go.AllowList", results[1].Name)
}

func TestSearchTemplates_NoMatch(t *testing.T) {
	resetSingleton()
	t.Setenv("HOME", t.TempDir())

	c := GetInstance()
	c.Templates = []Template{{Name: "Go"}}

	results := c.SearchTemplates("xyz")
	assert.Empty(t, results)
}

func TestFindTemplate_CaseInsensitive(t *testing.T) {
	resetSingleton()
	t.Setenv("HOME", t.TempDir())

	c := GetInstance()
	c.Templates = []Template{
		{Name: "Go", URL: "http://example.com"},
	}

	tmpl, err := c.FindTemplate("go")
	require.NoError(t, err)
	assert.Equal(t, "Go", tmpl.Name)
}

func TestSaveAndLoad(t *testing.T) {
	resetSingleton()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	c := GetInstance()
	c.Templates = []Template{{Name: "Go", URL: "http://example.com"}}
	require.NoError(t, c.Save())

	cachePath := filepath.Join(tmpDir, ".cache", "ignore-cli", "cache.json")
	_, err := os.Stat(cachePath)
	require.NoError(t, err, "cache file should exist")

	resetSingleton()
	c2 := GetInstance()
	require.Len(t, c2.Templates, 1)
	assert.Equal(t, "Go", c2.Templates[0].Name)
}
