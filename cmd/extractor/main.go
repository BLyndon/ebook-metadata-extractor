package main

import (
	"context"
	"ebook-metadata-extractor/config"
	"ebook-metadata-extractor/pkg/fileutil"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	cfg := config.LoadConfig()

	titles := fileutil.ReadTitles(cfg)

	for _, title := range titles {
		fmt.Printf("Processing: %v\n", title)
		err := extractMetaData(title, cfg)
		if err != nil {
			fmt.Printf("Error extracting metadata: %v\n", err)
		}
	}
}

func extractMetaData(title string, cfg config.Config) error {
	prompt, err := preparePrompt(cfg.PromptFile, title)
	if err != nil {
		fmt.Printf("Error preparing prompt: %v\n", err)
		return nil
	}

	client := getClient()
	ctx := context.Background()

	err = generateChatCompletion(ctx, client, title, prompt, cfg)
	if err != nil {
		fmt.Printf("Error in chat completion: %v\n", err)
	}
	return nil
}

func getClient() *openai.Client {
	return openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}

func preparePrompt(filePath, title string) (string, error) {
	prompt, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	promptString := strings.Replace(string(prompt), "{title}", title, -1)
	return promptString, nil
}

func generateChatCompletion(ctx context.Context, client *openai.Client, title string, prompt string, cfg config.Config) error {
	req := openai.ChatCompletionRequest{
		Model:     cfg.OpenAIModel,
		MaxTokens: cfg.MaxTokenCount,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("chat completion stream error: %v", err)
	}
	defer stream.Close()

	return handleStreamResponse(stream, title, cfg)
}

func handleStreamResponse(stream *openai.ChatCompletionStream, title string, cfg config.Config) error {
	targetFile := cfg.TargetDir + title + cfg.JsonExtension

	file, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("stream error: %v", err)
		}

		_, writeErr := file.WriteString(response.Choices[0].Delta.Content)
		if writeErr != nil {
			return fmt.Errorf("error writing to file: %v", writeErr)
		}
	}

	return nil
}
