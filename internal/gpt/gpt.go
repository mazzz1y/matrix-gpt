package gpt

import (
	"time"

	"github.com/sashabaranov/go-openai"
)

type Gpt struct {
	client      *openai.Client
	model       string
	gptTimeout  time.Duration
	maxAttempts int
}

// New initializes a Gpt instance with the provided configurations.
func New(token, gptModel string, historyLimit, gptTimeout, maxAttempts int) *Gpt {
	return &Gpt{
		client:      openai.NewClient(token),
		model:       gptModel,
		gptTimeout:  time.Duration(gptTimeout) * time.Second,
		maxAttempts: maxAttempts,
	}
}

// GetModel returns the GPT model string.
func (g *Gpt) GetModel() string {
	return g.model
}

// GetTimeout returns the timeout value for the GPT client.
func (g *Gpt) GetTimeout() time.Duration {
	return g.gptTimeout
}
