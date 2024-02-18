package fileutil

import (
	"ebook-metadata-extractor/config"
	"fmt"
	"os"
)

func ReadTitles(cfg config.Config) []string {
	files, err := os.ReadDir(cfg.SourceDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return nil
	}

	var titles []string
	for _, file := range files {
		if file.Name() != ".DS_Store" {
			titles = append(titles, file.Name())
		}
	}
	return titles
}
