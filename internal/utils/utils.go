package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	cacheFile        = "cache.json"
	cacheDirName     = "ignore-cli"
	cacheTemplateDir = "templates"
)

func GetCacheFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, ".cache", cacheDirName)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return filepath.Join(cacheDir, cacheFile), nil
}

func GetTemplateDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	templateDir := filepath.Join(homeDir, ".cache", cacheDirName, cacheTemplateDir)
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create templates directory: %w", err)
	}

	return templateDir, nil
}
