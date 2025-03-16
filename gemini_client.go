package main

import (
	"context"
	"log"
	"os"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	Client *genai.Client
	Model *genai.EmbeddingModel
}


func NewGeminiClient() *GeminiClient{
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	em := client.EmbeddingModel("gemini-embedding-exp-03-07")
	return &GeminiClient{
		Client: client,
		Model: em,
	}
}
func (geminiClient *GeminiClient) embed(key string) ([]float32,error){
	ctx := context.Background()
	res, err := geminiClient.Model.EmbedContent(ctx, genai.Text(key))
	if err != nil {
		return nil,err
	}
	return res.Embedding.Values,nil
}
