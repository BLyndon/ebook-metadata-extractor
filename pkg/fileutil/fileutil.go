package fileutil

import (
	"ebook-metadata-extractor/config"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func ReadTitles(cfg config.Config) []string {
	sourceFiles := getSourceFiles(cfg)

	var titles []string
	for _, filename := range sourceFiles {
		if filename != ".DS_Store" {
			titles = append(titles, filename)
		}
	}
	return titles
}

func HandleStreamResponse(stream *openai.ChatCompletionStream, title string, cfg config.Config) (string, error) {
	targetFile := filepath.Join(cfg.TargetDir, title+cfg.TargetFileExtension)
	jsonResponse := ""

	file, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("stream error: %v", err)
		}

		jsonPart := response.Choices[0].Delta.Content

		jsonResponse += jsonPart
		if cfg.PeristMetadata {
			_, err = file.WriteString(jsonPart)
			if err != nil {
				log.Fatalf("error writing to file: %v", err)
			}
		}
	}

	log.Printf("JSON response: %v\n", jsonResponse)

	return jsonResponse, nil
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
	var filesWithoutExtensions []string
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
