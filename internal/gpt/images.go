package gpt

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"
)

// CreateImage makes a request to get a DALL-E image URL.
func (g *Gpt) CreateImage(ctx context.Context, style, prompt string) (string, error) {
	var res openai.ImageResponse
	var err error

	for i := 0; i < g.maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(ctx, g.gptTimeout)
		defer cancel()

		res, err = g.client.CreateImage(
			ctx,
			openai.ImageRequest{
				Model:          openai.CreateImageModelDallE3,
				Style:          style,
				Prompt:         prompt,
				Size:           "1024x1024",
				ResponseFormat: openai.CreateImageResponseFormatURL,
			},
		)

		if ctx.Err() == context.Canceled {
			return "", ctx.Err()
		} else if !isServiceUnavailableError(err) {
			break
		}

		sleepBeforeRetry(i)
	}

	if len(res.Data) < 1 {
		return "", errors.New("empty response")
	}

	return res.Data[0].URL, err
}
