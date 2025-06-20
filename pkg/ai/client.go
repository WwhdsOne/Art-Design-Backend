package ai

import (
	"net/http"
)

type AIModelClient struct {
	client *http.Client
}

func NewAIModelClient(client *http.Client) *AIModelClient {
	return &AIModelClient{
		client: client,
	}
}
