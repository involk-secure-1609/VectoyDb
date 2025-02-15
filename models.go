package main

type EmbedResponsePayload struct {
	Model           string    `json:"model"`
	Embeddings      [][]float32 `json:"embeddings"`
	TotalDuration   int64     `json:"total_duration"`
	LoadDuration    int64     `json:"load_duration"`
	PromptEvalCount int       `json:"prompt_eval_count"`
}
// RequestPayload represents the payload sent in the HTTP request.
type EmbedRequestPayload struct {
	Model      string   `json:"model"`
	Input     []string   `json:"input"`
}
