package client

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// EmbeddingModel interface abstracts the genai.EmbeddingModel
type EmbeddingModel interface {
	EmbedContent(ctx context.Context,parts ...genai.Part) (*genai.EmbedContentResponse, error)
}

// ClientFactory creates genai clients
type ClientFactory interface {
	NewClient(ctx context.Context, opts ...option.ClientOption) (*genai.Client, error)
	GetEmbeddingModel(client *genai.Client, modelName string) EmbeddingModel
}

// DefaultClientFactory is the default implementation of ClientFactory
type DefaultClientFactory struct{}

func (factory *DefaultClientFactory) NewClient(ctx context.Context, opts ...option.ClientOption) (*genai.Client, error) {
	return genai.NewClient(ctx, opts...)
}

func (factory *DefaultClientFactory) GetEmbeddingModel(client *genai.Client, modelName string) EmbeddingModel {
	return client.EmbeddingModel(modelName)
}

// GeminiClient provides access to Gemini API
type GeminiClient struct {
	Client *genai.Client
	Model  EmbeddingModel
}

// NewGeminiClient creates a new client for Gemini API
func NewGeminiClient() *GeminiClient {
	return NewGeminiClientWithFactory(&DefaultClientFactory{})
}

// NewGeminiClientWithFactory creates a new client using the provided factory
func NewGeminiClientWithFactory(factory ClientFactory) *GeminiClient {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read environment variables
	geminiApiKey := os.Getenv("GEMINI_API_KEY")
	client, err := factory.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		log.Fatal(err)
	}

	model := factory.GetEmbeddingModel(client, "gemini-embedding-exp-03-07")
	return &GeminiClient{
		Client: client,
		Model:  model,
	}
}

// Embed creates an embedding for the given key
func (geminiClient *GeminiClient) Embed(key string) ([]float32, error) {
	ctx := context.Background()
	res, err := geminiClient.Model.EmbedContent(ctx, genai.Text(key))
	if err != nil {
		return nil, err
	}
	return res.Embedding.Values, nil
}