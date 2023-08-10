package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"time"
)

// CreateImage makes a request to get a DALL-E image URL.
func (g *Gpt) CreateImage(prompt string) (openai.ImageResponse, error) {
	ctx, cancel := context.WithTimeout(g.ctx, time.Duration(g.gptTimeout)*time.Second)
	defer cancel()

	return g.client.CreateImage(
		ctx,
		openai.ImageRequest{
			Prompt:         prompt,
			Size:           "1024x1024",
			ResponseFormat: openai.CreateImageResponseFormatURL,
		},
	)
}