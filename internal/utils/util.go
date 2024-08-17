package utils

import (
	"os"
	"path"
)

const (
	cacheFile        = "cache.json"
	cacheDir         = "ignore-cli"
	cacheTemplateDir = "templates"
)

func GetCacheFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	cacheFilePath := path.Join(homeDir, ".cache", cacheDir, cacheFile)

	// create cache directory if it doesn't exist
	cacheDir := path.Join(homeDir, ".cache", cacheDir)
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, os.ModePerm)
	}

	return cacheFilePath
}

func GetTemplateDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	templateDir := path.Join(homeDir, ".cache", cacheDir, cacheTemplateDir)

	// create templates directory if it doesn't exist
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		os.MkdirAll(templateDir, os.ModePerm)
	}

	return templateDir
}
