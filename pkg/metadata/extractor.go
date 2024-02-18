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

	client := getClient()
	ctx := context.Background()

	generateChatCompletion(ctx, client, title, prompt, cfg)
}

func getClient() *openai.Client {
	return openai.NewClient(os.Getenv("OPENAI_API_KEY"))
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

func generateChatCompletion(ctx context.Context, client *openai.Client, title string, prompt string, cfg config.Config) error {
	req := configureChatCompletionRequest(prompt, cfg)

	stream, err := client.CreateChatCompletionStream(ctx, *req)
	if err != nil {
		return fmt.Errorf("chat completion stream error: %v", err)
	}
	defer stream.Close()

	return fileutil.HandleStreamResponse(stream, title, cfg)
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
