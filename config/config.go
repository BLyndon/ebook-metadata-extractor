package config

import (
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type Config struct {
	OpenAIModel         string
	MaxTokenCount       int
	PromptFile          string
	TargetFileExtension string
	SourceDir           string
	TargetDir           string
	PeristMetadata      bool
	Port                int
	OpenApiKey          string
}

func LoadConfig() Config {
	return Config{
		OpenAIModel:         openai.GPT3Dot5Turbo,
		MaxTokenCount:       1000,
		PromptFile:          "./assets/prompt.txt",
		TargetFileExtension: ".json",
		SourceDir:           "./data/sourceDir",
		TargetDir:           "./data/targetDir",
		PeristMetadata:      true,
		Port:                8080,
		OpenApiKey:          os.Getenv("OPENAI_API_KEY"),
	}
}
