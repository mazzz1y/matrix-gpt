package gpt

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// CreateTranscription retrieves a transcription from audio file.
func (g *Gpt) CreateTranscription(ctx context.Context, fname string) (string, error) {
	var res openai.AudioResponse
	var err error

	for i := 0; i < g.maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(ctx, g.gptTimeout)
		defer cancel()

		res, err = g.client.CreateTranscription(
			context.Background(),
			openai.AudioRequest{
				Model:    openai.Whisper1,
				FilePath: fname,
			},
		)

		if ctx.Err() == context.Canceled {
			return "", ctx.Err()
		} else if !isServiceUnavailableError(err) {
			break
		}

		sleepBeforeRetry(i)
	}

	return res.Text, err
}
