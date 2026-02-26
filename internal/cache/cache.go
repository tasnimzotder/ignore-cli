package cache

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/tasnimzotder/ignore-cli/internal/utils"
)

type Cache struct {
	mu         sync.Mutex
	Templates  []Template `json:"templates"`
	LastUpdate time.Time  `json:"last_update"`
}

type Template struct {
	Name          string    `json:"name"`
	URL           string    `json:"url"`
	CacheFilePath string    `json:"cache_file_path"`
	LastUpdate    time.Time `json:"last_update"`
}

type gitignoreAPIResponse struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

var (
	instance         *Cache
	once             sync.Once
	gitignoreListURL = "https://api.github.com/gitignore/templates"
	gitignoreAPIURL  = "https://api.github.com/gitignore/templates/%s"
	HTTPClient       = &http.Client{
		Timeout: 30 * time.Second,
	}
)

func GetInstance() *Cache {
	once.Do(func() {
		instance = &Cache{}
		instance.load()
	})
	return instance
}

func (c *Cache) load() error {
	cacheFilePath, err := utils.GetCacheFilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, c)
}

func (c *Cache) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheFilePath, err := utils.GetCacheFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFilePath, data, 0o644)
}

func (c *Cache) ListTemplates() ([]Template, error) {
	if len(c.Templates) == 0 || time.Since(c.LastUpdate) > 7*24*time.Hour {
		if err := c.fetchTemplateList(); err != nil {
			return nil, err
		}
	}
	return c.Templates, nil
}

func (c *Cache) SearchTemplates(query string) []Template {
	query = strings.ToLower(query)
	var results []Template
	for _, t := range c.Templates {
		if strings.Contains(strings.ToLower(t.Name), query) {
			results = append(results, t)
		}
	}
	return results
}

func (c *Cache) FindTemplate(name string) (*Template, error) {
	for i := range c.Templates {
		if strings.EqualFold(c.Templates[i].Name, name) {
			return &c.Templates[i], nil
		}
	}

	t := Template{
		Name: name,
		URL:  fmt.Sprintf(gitignoreAPIURL, name),
	}
	if err := t.Update(); err != nil {
		return nil, fmt.Errorf("template %q not found: %w", name, err)
	}
	c.UpdateTemplate(t)
	if err := c.Save(); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *Cache) UpdateTemplate(t Template) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, existing := range c.Templates {
		if existing.Name == t.Name {
			c.Templates[i] = t
			return
		}
	}
	c.Templates = append(c.Templates, t)
}

func (c *Cache) fetchTemplateList() error {
	resp, err := HTTPClient.Get(gitignoreListURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch templates: %s", resp.Status)
	}

	var names []string
	if err := json.NewDecoder(resp.Body).Decode(&names); err != nil {
		return err
	}

	now := time.Now()
	for _, name := range names {
		c.UpdateTemplate(Template{
			Name:       name,
			URL:        fmt.Sprintf(gitignoreAPIURL, name),
			LastUpdate: now,
		})
	}
	c.LastUpdate = now
	return c.Save()
}

func (t *Template) NeedsUpdate() bool {
	return t.CacheFilePath == "" || time.Since(t.LastUpdate) > 24*time.Hour
}

func (t *Template) Update() error {
	resp, err := HTTPClient.Get(t.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch template: %s", resp.Status)
	}

	var apiResp gitignoreAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode API response: %w", err)
	}

	cacheDir, err := utils.GetTemplateDir()
	if err != nil {
		return err
	}

	t.CacheFilePath = filepath.Join(cacheDir, t.Name+".gitignore")
	if err := os.WriteFile(t.CacheFilePath, []byte(apiResp.Source), 0o644); err != nil {
		return err
	}

	t.LastUpdate = time.Now()
	return nil
}

func (t *Template) Content() (string, error) {
	if t.NeedsUpdate() {
		if err := t.Update(); err != nil {
			return "", err
		}
	}

	content, err := os.ReadFile(t.CacheFilePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
