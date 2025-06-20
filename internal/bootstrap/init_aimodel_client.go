package bootstrap

import (
	"Art-Design-Backend/pkg/ai"
	"net/http"
	"time"
)

func InitAIModelClient() *ai.AIModelClient {
	h := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
	return ai.NewAIModelClient(h)
}
