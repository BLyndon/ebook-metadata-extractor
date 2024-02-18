package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const (
	OpenAIModel   = openai.GPT3Dot5Turbo
	MaxTokenCount = 1000
	PromptFile    = "./prompt.txt"
	JsonExtension = ".json"
	SourceDir     = "./sourceDir"
	TargetDir     = "./targetDir/"
)

func main() {
	titles := getTitles()

	for _, title := range titles {
		err := extractMetaData(title)
		if err != nil {
			fmt.Printf("Error extracting metadata: %v\n", err)
		}
	}
}

func getTitles() []string {
	files, err := os.ReadDir(SourceDir)
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

func extractMetaData(title string) error {
	prompt, err := preparePrompt(PromptFile, title)
	if err != nil {
		fmt.Printf("Error preparing prompt: %v\n", err)
		return nil
	}

	fmt.Printf("Prompt: %v\n", prompt)

	client := getClient()
	ctx := context.Background()

	err = generateChatCompletion(ctx, client, title, prompt)
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

func generateChatCompletion(ctx context.Context, client *openai.Client, title string, prompt string) error {
	req := openai.ChatCompletionRequest{
		Model:     OpenAIModel,
		MaxTokens: MaxTokenCount,
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

	return handleStreamResponse(stream, title)
}

func handleStreamResponse(stream *openai.ChatCompletionStream, title string) error {
	targetFile := TargetDir + title + JsonExtension

	file, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Extraction finished")
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