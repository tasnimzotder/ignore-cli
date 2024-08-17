package search

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/tasnimzotder/ignore-cli/internal/cache"
	"github.com/tasnimzotder/ignore-cli/internal/utils"
)

const (
	gitignoreListURL    = "https://api.github.com/gitignore/templates"
	gitignoreContentURL = "https://raw.githubusercontent.com/github/gitignore/main/%s.gitignore"
)

func Templates(query string) ([]string, error) {
	_cache, err := cache.Get()
	if err != nil {
		return nil, err
	}

	if len(_cache.Templates) == 0 {
		err = UpdateTemplateList(_cache)
		if err != nil {
			return nil, err
		}
	}

	var result []string
	for idx := range _cache.Templates {
		template := _cache.Templates[idx]
		if strings.Contains(strings.ToLower(template.Name), strings.ToLower(query)) {
			result = append(result, template.Name)
		}
	}

	return result, nil
}

func AllTemplates() ([]string, error) {
	cache, err := cache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	if len(cache.Templates) == 0 {
		log.Println("Cache is empty. Updating cache...")
		err = UpdateTemplateList(cache)
		if err != nil {
			return nil, err
		}
	}

	templates := make([]string, 0, len(cache.Templates))
	for idx := range cache.Templates {
		templates = append(templates, cache.Templates[idx].Name)
	}

	return templates, nil
}

func UpdateTemplateList(_cache *cache.Cache) error {
	resp, err := http.Get(gitignoreListURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch .gitignore templates: %s", resp.Status)
	}

	var templates []string
	err = json.NewDecoder(resp.Body).Decode(&templates)
	if err != nil {
		return err
	}

	_cache.Templates = make([]cache.Template, len(templates))

	templatesDir := utils.GetTemplateDir()

	var wg sync.WaitGroup

	for idx, template := range templates {
		_cache.Templates[idx].Name = template
		_cache.Templates[idx].URL = fmt.Sprintf(gitignoreContentURL, template)

		// create the cache file path
		// _cache.Templates[idx].CacheFilePath = fmt.Sprintf("templates/%s.gitignore", template)
		_cache.Templates[idx].CacheFilePath = fmt.Sprintf("%s/%s.gitignore", templatesDir, template)

		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// create the templates directory
			err = os.MkdirAll("templates", os.ModePerm)
			if err != nil {
				return
			}

			// create the file
			file, err := os.Create(_cache.Templates[idx].CacheFilePath)
			if err != nil {
				return
			}
			defer file.Close()

			// get the content
			resp, err := http.Get(_cache.Templates[idx].URL)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			// write the content to the file
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				return
			}
		}(idx)
	}

	wg.Wait()

	return _cache.Save()
}

func GetTemplateContent(name string) (string, error) {
	_cache, err := cache.Get()
	if err != nil {
		return "", err
	}

	for _, template := range _cache.Templates {
		if strings.EqualFold(template.Name, name) {
			filePath := template.CacheFilePath
			// check if the file exists
			// if it doesn't exist, fetch the content from the API
			info, err := os.Stat(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					content, err := fetchAndSaveTemplateContent(template, _cache)
					if err != nil {
						return "", err
					}

					return content, nil
				}

				return "", err
			}

			// check if the file is empty
			if info.Size() == 0 {
				content, err := fetchAndSaveTemplateContent(template, _cache)
				if err != nil {
					return "", err
				}

				return content, nil
			}

			// read the content from the file
			data, err := os.ReadFile(filePath)
			if err != nil {
				return "", err
			}

			return string(data), nil
		}
	}

	return "", nil
}

func fetchAndSaveTemplateContent(template cache.Template, cache *cache.Cache) (string, error) {
	resp, err := http.Get(template.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch .gitignore template: %s", resp.Status)
	}

	// create the file
	file, err := os.Create(template.CacheFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// write the content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	// update the cache
	for idx := range cache.Templates {
		if cache.Templates[idx].Name == template.Name {
			cache.Templates[idx] = template
		}
	}

	err = cache.Save()
	if err != nil {
		return "", err
	}

	// read the content from the file
	data, err := os.ReadFile(template.CacheFilePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
