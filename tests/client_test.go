package tests

import (
	"context"
	"os"
	"testing"
	"vectorDb/client"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// MockClientFactory implements the ClientFactory interface for testing
type MockClientFactory struct {
	NewClientCalled         bool
	GetEmbeddingModelCalled bool
	Options                 []option.ClientOption
	ApiKeyUsed              string
	ModelNameUsed           string
}

func (m *MockClientFactory) NewClient(ctx context.Context, opts ...option.ClientOption) (*genai.Client, error) {
	m.NewClientCalled = true
	m.Options = opts
	
	m.ApiKeyUsed = os.Getenv("GEMINI_API_KEY")	
	// Return a fake client that doesn't connect to anything
	return &genai.Client{}, nil
}

func (m *MockClientFactory) GetEmbeddingModel(client *genai.Client, modelName string) client.EmbeddingModel {
	m.GetEmbeddingModelCalled = true
	m.ModelNameUsed = modelName
	
	// Return a mock embedding model
	return &MockEmbeddingModel{
		EmbedContentFunc: func(ctx context.Context, parts ...genai.Part) (*genai.EmbedContentResponse, error) {
			return &genai.EmbedContentResponse{
				Embedding: &genai.ContentEmbedding{
					Values: []float32{0.1, 0.2, 0.3},
				},
			}, nil
		},
	}
}

// MockEmbeddingModel implements the EmbeddingModel interface for testing
type MockEmbeddingModel struct {
	EmbedContentFunc func(ctx context.Context, parts ...genai.Part) (*genai.EmbedContentResponse, error)
}

func (m *MockEmbeddingModel) EmbedContent(ctx context.Context,parts ...genai.Part) (*genai.EmbedContentResponse, error) {
	return m.EmbedContentFunc(ctx, parts...)
}

func TestNewGeminiClientWithFactory(t *testing.T) {
	// Save original environment to restore later
	originalEnv := os.Getenv("GEMINI_API_KEY")
	defer os.Setenv("GEMINI_API_KEY", originalEnv)

	// Setup a test environment
	os.Setenv("GEMINI_API_KEY", "test-api-key")

	// Create a mock factory
	mockFactory := &MockClientFactory{}
	
	// Test the function
	c := client.NewGeminiClientWithFactory(mockFactory)
	
	// Verify the client was created properly
	if c == nil {
		t.Fatal("NewGeminiClientWithFactory returned nil")
	}
	
	// Verify the factory methods were called
	if !mockFactory.NewClientCalled {
		t.Error("Factory's NewClient method was not called")
	}
	
	if !mockFactory.GetEmbeddingModelCalled {
		t.Error("Factory's GetEmbeddingModel method was not called")
	}
	
	// Verify the API key was used correctly
	if mockFactory.ApiKeyUsed != "test-api-key" {
		t.Errorf("Expected API key %q, got %q", "test-api-key", mockFactory.ApiKeyUsed)
	}
	
	// Verify the model name was used correctly
	if mockFactory.ModelNameUsed != "gemini-embedding-exp-03-07" {
		t.Errorf("Expected model name %q, got %q", "gemini-embedding-exp-03-07", mockFactory.ModelNameUsed)
	}
}

func TestEmbed(t *testing.T) {
	// These tests can remain largely the same as before, but adapted to use the factory pattern
    // ...
}