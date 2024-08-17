package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/tasnimzotder/ignore-cli/internal/utils"
)

type Cache struct {
	Templates  []Template `json:"templates"`
	LastUpdate time.Time  `json:"last_update"`
}

type Template struct {
	Name          string    `json:"name"`
	URL           string    `json:"url"`
	CacheFilePath string    `json:"cache_file_path"`
	LastUpdate    time.Time `json:"last_update"`
}

func Get() (*Cache, error) {
	// homeDir, err := os.UserHomeDir()
	// homeDir, err := os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }

	cacheFilePath := utils.GetCacheFilePath()
	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Cache{Templates: []Template{}}, nil
		}
		return nil, err
	}

	var cache Cache
	err = json.Unmarshal(data, &cache)
	return &cache, err
}

func (c *Cache) Save() error {
	// homeDir, err := os.UserHomeDir()
	// homeDir, err := os.Getwd()
	// if err != nil {
	// 	return err
	// }

	cacheFilePath := utils.GetCacheFilePath()
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFilePath, data, 0644)
}

func (c *Cache) UpdateTemplate(t Template) {
	for i, existingTemplate := range c.Templates {
		if existingTemplate.Name == t.Name {
			c.Templates[i] = t
			return
		}
	}
	c.Templates = append(c.Templates, t)
}

func (t *Template) NeedsUpdate() bool {
	return t.CacheFilePath == "" || time.Since(t.LastUpdate) > 24*time.Hour
}

func (t *Template) Update() error {
	resp, err := http.Get(t.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch gitignore template: %s", resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// homeDir, err := os.UserHomeDir()
	// homeDir, err := os.Getwd()
	// if err != nil {
	// 	return err
	// }

	// cacheDir := filepath.Join(homeDir, "templates")
	cacheDir := utils.GetTemplateDir()
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	t.CacheFilePath = filepath.Join(cacheDir, t.Name+".gitignore")
	if err := os.WriteFile(t.CacheFilePath, content, 0644); err != nil {
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
