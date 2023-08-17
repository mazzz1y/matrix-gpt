package gpt

import (
	"context"
	"errors"
	"time"

	"github.com/sashabaranov/go-openai"
)

// CreateImage makes a request to get a DALL-E image URL.
func (g *Gpt) CreateImage(ctx context.Context, prompt string) (openai.ImageResponse, error) {
	var response openai.ImageResponse
	var err error

	ctx, cancel := context.WithTimeout(ctx, g.gptTimeout)
	defer cancel()

	for i := 0; i < g.maxAttempts; i++ {
		response, err = g.client.CreateImage(
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
		if err == nil && len(response.Data) < 1 {
			err = errors.New("empty response")
		} else if err == nil {
			return response, nil
		}
		time.Sleep(time.Duration(i*2) * time.Second)
	}

	return openai.ImageResponse{}, err
}
