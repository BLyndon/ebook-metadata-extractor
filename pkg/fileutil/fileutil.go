package fileutil

import (
	"ebook-metadata-extractor/config"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func WriteToFile(metadataString string, title string, cfg config.Config) error {
	filePath := filepath.Join(cfg.TargetDir, title+cfg.TargetFileExtension)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("error creating file: %v\n", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(metadataString)
	if err != nil {
		log.Printf("error writing to file: %v\n", err)
		return err
	}
	log.Printf("Metadata written to: %v\n", filePath)

	return nil
}

func GetSourceTargetTitleDelta(cfg config.Config) []string {
	sourceFiles := getSourceFiles(cfg)

	var titles []string
	for _, filename := range sourceFiles {
		if filename != ".DS_Store" {
			titles = append(titles, filename)
		}
	}
	return titles
}

func ReadMetadataIfExists(title string, cfg config.Config) (string, error) {
	fileName := filepath.Join(cfg.TargetDir, title+cfg.TargetFileExtension)
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	fmt.Printf("Metadata found for title: %v\n", title)
	return string(content), nil
}

func getSourceFiles(cfg config.Config) []string {
	sourceFiles := readAllFileNamesIn(cfg.SourceDir)
	targetFiles := readAllFileNamesIn(cfg.TargetDir)

	for _, targetFile := range targetFiles {
		for i, sourceFile := range sourceFiles {
			if sourceFile == targetFile {
				sourceFiles = append(sourceFiles[:i], sourceFiles[i+1:]...)
			}
		}
	}

	return sourceFiles
}

func readAllFileNamesIn(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("error reading directory: %v\n", err)
		return nil
	}
	return removeFileExtensions(files)
}

func removeFileExtensions(files []fs.DirEntry) []string {
	filesWithoutExtensions := make([]string, 0)
	for _, file := range files {
		if file.Name() != ".DS_Store" {
			baseFileName := getBaseFileName(file.Name())
			filesWithoutExtensions = append(filesWithoutExtensions, baseFileName)
		}
	}
	return filesWithoutExtensions
}

func getBaseFileName(fileName string) string {
	baseFileName := strings.Replace(fileName, ".pdf", "", -1)
	baseFileName = strings.Replace(baseFileName, ".epub", "", -1)
	baseFileName = strings.Replace(baseFileName, ".mobi", "", -1)
	baseFileName = strings.Replace(baseFileName, ".json", "", -1)
	return baseFileName
}
