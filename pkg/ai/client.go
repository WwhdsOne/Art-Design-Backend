package ai

import (
	"net/http"
	"time"
)

type AIModelClient struct {
	client *http.Client
}

func NewAIModelClient() *AIModelClient {
	return &AIModelClient{
		client: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}
