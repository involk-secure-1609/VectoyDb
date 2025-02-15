package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OllamaClient struct {
	BaseUrl    string
	HttpClient *http.Client
	model      string
}

// var (
// 	once                 sync.Once
// 	ollamaClientInstance *OllamaClient
// )

func NewOllamaClient(baseUrl string, model string) *OllamaClient {
	// once.Do(func() {

		httpClient := http.Client{
			Timeout: 0,
		}
		ollamaClientInstance := &OllamaClient{
			HttpClient: &httpClient,
			BaseUrl:    baseUrl,
			model:      model,
		}

	// })

	return ollamaClientInstance
}

func (client *OllamaClient) start() error {
	isServerRunning := checkIfOllamaServerRunning(client.BaseUrl)
	if !isServerRunning {
		return fmt.Errorf("the ollama server has not started ")
	}
	return nil
}

func (client *OllamaClient) embed(inputs []string) (*http.Response, error) {

	requestPayload := EmbedRequestPayload{
		Model: client.model,
		Input: inputs,
	}

	jsonRequest, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	log.Println("JSON payload marshaled successfully")

	// Log the request payload
	log.Printf("Request payload: %s", string(jsonRequest))
	log.Printf("Creating new HTTP POST request to %s", client.BaseUrl)
	req, err := http.NewRequest("POST", "http://"+client.BaseUrl+"/api/embed", bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Printf("HTTP request creation error: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	log.Println("Sending HTTP request")
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		log.Printf("HTTP request error: %v", err)
		return nil, err
	}

	log.Printf("Received HTTP response with status code: %d", resp.StatusCode)
	return resp, nil
}
