package main

import (
	"os"
	"path/filepath"

	"github.com/wizenheimer/cascade/internal/config"
	"gopkg.in/yaml.v2"
)

// compose creates a yaml file from the given configuration
func compose(cfg config.Config, fileName, outputFolder string, createFolder bool) error {
	yamlData, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	if fileName == "" {
		fileName = config.GetEnv("FILE_NAME", config.FILE_NAME)
	}

	if outputFolder == "" {
		outputFolder = config.GetEnv("OUTPUT_FOLDER", config.OUTPUT_FOLDER)
	}

	fullPath := filepath.Join(outputFolder, fileName)

	if createFolder {
		err = os.MkdirAll(outputFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(fullPath, yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}
