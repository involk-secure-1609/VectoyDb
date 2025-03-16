package client

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	Client *genai.Client
	Model  *genai.EmbeddingModel
}

func NewGeminiClient() *GeminiClient {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read environment variables
	geminiApiKey := os.Getenv("GEMINI_API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		log.Fatal(err)
	}
	em := client.EmbeddingModel("gemini-embedding-exp-03-07")
	return &GeminiClient{
		Client: client,
		Model:  em,
	}
}
func (geminiClient *GeminiClient) Embed(key string) ([]float32, error) {
	ctx := context.Background()
	res, err := geminiClient.Model.EmbedContent(ctx, genai.Text(key))
	if err != nil {
		return nil, err
	}
	return res.Embedding.Values, nil
}
