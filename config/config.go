package config

import (
	"fmt"
	"os"
	"strconv"

	openai "github.com/sashabaranov/go-openai"
)

type Config struct {
	OpenAIModel         string
	MaxTokenCount       int
	PromptFile          string
	TargetFileExtension string
	SourceDir           string
	TargetDir           string
	PersistMetadata     bool
	Port                int
	OpenApiKey          string
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	fmt.Printf("Using default value for %s\n", key)
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	fmt.Printf("Using default value for %s\n", key)
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	fmt.Printf("Using default value for %s\n", key)
	return fallback
}

func LoadConfig() Config {
	return Config{
		OpenAIModel:         getEnv("OPENAI_MODEL", openai.GPT3Dot5Turbo),
		MaxTokenCount:       getEnvAsInt("MAX_TOKEN_COUNT", 1000),
		PromptFile:          getEnv("PROMPT_FILE", "./assets/prompt.txt"),
		TargetFileExtension: getEnv("TARGET_FILE_EXTENSION", ".json"),
		SourceDir:           getEnv("SOURCE_DIR", "./data/sourceDir"),
		TargetDir:           getEnv("TARGET_DIR", "./data/targetDir"),
		PersistMetadata:     getEnvAsBool("PERSIST_METADATA", true),
		Port:                getEnvAsInt("PORT", 8080),
		OpenApiKey:          os.Getenv("OPENAI_API_KEY"),
	}
}
