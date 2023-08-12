package gpt

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// CreateImage makes a request to get a DALL-E image URL.
func (g *Gpt) CreateImage(ctx context.Context, prompt string) (openai.ImageResponse, error) {
	var err error
	for i := 0; i < g.maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(ctx, g.gptTimeout)
		defer cancel()

		response, err := g.client.CreateImage(
			ctx,
			openai.ImageRequest{
				Prompt:         prompt,
				Size:           "1024x1024",
				ResponseFormat: openai.CreateImageResponseFormatURL,
			},
		)
		if ctx.Err() == context.Canceled {
			return openai.ImageResponse{}, ctx.Err()
		}
		if err == nil {
			return response, nil
		}
	}

	return openai.ImageResponse{}, err
}
