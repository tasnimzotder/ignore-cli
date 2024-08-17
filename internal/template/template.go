package template

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tasnimzotder/ignore-cli/internal/cache"
)

const (
	gitignoreAPIURL  = "https://api.github.com/gitignore/templates/%s"
	gitignoreListURL = "https://api.github.com/gitignore/templates"
)

func Get(name string) (*cache.Template, error) {
	c, err := cache.Get()
	if err != nil {
		return nil, err
	}

	for _, t := range c.Templates {
		if strings.EqualFold(t.Name, name) {
			if t.NeedsUpdate() {
				if err := t.Update(); err != nil {
					return nil, err
				}
				c.UpdateTemplate(t)
				if err := c.Save(); err != nil {
					return nil, err
				}
			}
			return &t, nil
		}
	}

	// Template not found in cache, fetch it
	t := cache.Template{
		Name: name,
		URL:  fmt.Sprintf(gitignoreAPIURL, name),
	}
	if err := t.Update(); err != nil {
		return nil, err
	}
	c.UpdateTemplate(t)
	if err := c.Save(); err != nil {
		return nil, err
	}

	return &t, nil
}

func Search(query string) ([]cache.Template, error) {
	templates, err := List()
	if err != nil {
		return nil, err
	}

	var results []cache.Template
	for _, t := range templates {
		if strings.Contains(strings.ToLower(t.Name), strings.ToLower(query)) {
			results = append(results, t)
		}
	}

	return results, nil
}

func List() ([]cache.Template, error) {
	c, err := cache.Get()
	if err != nil {
		return nil, err
	}

	if len(c.Templates) == 0 || time.Since(c.LastUpdate) > 24*time.Hour {
		if err := updateTemplateList(c); err != nil {
			return nil, err
		}
	}

	return c.Templates, nil
}

func updateTemplateList(c *cache.Cache) error {
	resp, err := http.Get(gitignoreListURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch gitignore templates: %s", resp.Status)
	}

	var templateNames []string
	err = json.NewDecoder(resp.Body).Decode(&templateNames)
	if err != nil {
		return err
	}

	for _, name := range templateNames {
		t := cache.Template{
			Name: name,
			URL:  fmt.Sprintf(gitignoreAPIURL, name),
		}
		c.UpdateTemplate(t)
	}

	c.LastUpdate = time.Now()
	return c.Save()
}
