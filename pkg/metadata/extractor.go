package metadata

import (
	"context"
	"ebook-metadata-extractor/config"
	"ebook-metadata-extractor/pkg/fileutil"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func ExtractMetaData(title string, cfg config.Config) {
	fmt.Printf("Processing: %v\n", title)

	prompt := preparePrompt(cfg.PromptFile, title)
	req := configureChatCompletionRequest(prompt, cfg)

	stream := generateChatCompletion(req, cfg)
	fileutil.HandleStreamResponse(stream, title, cfg)
}

func preparePrompt(filePath, title string) string {
	prompt, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("error reading prompt file: %v\n", err)
	}

	promptString := strings.Replace(string(prompt), "{title}", title, -1)

	if err != nil {
		fmt.Printf("error evaluating prompt template: %v\n", err)
	}
	return promptString
}

func generateChatCompletion(req *openai.ChatCompletionRequest, cfg config.Config) *openai.ChatCompletionStream {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	stream, err := client.CreateChatCompletionStream(context.Background(), *req)
	if err != nil {
		fmt.Printf("chat completion stream error: %v", err)
	}

	return stream
}

func configureChatCompletionRequest(prompt string, cfg config.Config) *openai.ChatCompletionRequest {
	return &openai.ChatCompletionRequest{
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
}
