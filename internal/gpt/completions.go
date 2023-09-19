package gpt

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"
)

// CreateCompletion retrieves a completion from GPT using the given user's message.
func (g *Gpt) CreateCompletion(ctx context.Context, history []openai.ChatCompletionMessage, userMsg string) ([]openai.ChatCompletionMessage, error) {
	// Append the user's message to the existing history.
	messageHistory := append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMsg,
	})

	res, err := g.complReqWithTimeout(ctx, messageHistory)
	if err != nil {
		return []openai.ChatCompletionMessage{}, err
	}

	messageHistory = append(messageHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: res,
	})

	return messageHistory, err
}

// complReqWithTimeout makes a request to get a GPT completion with a specified timeout.
func (g *Gpt) complReqWithTimeout(ctx context.Context, msg []openai.ChatCompletionMessage) (string, error) {
	var res openai.ChatCompletionResponse
	var err error

	for i := 0; i < g.maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(ctx, g.gptTimeout)
		defer cancel()

		res, err = g.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    g.model,
				Messages: msg,
			},
		)

		if ctx.Err() == context.Canceled {
			return "", ctx.Err()
		} else if isTokenExceededError(err) {
			msg = trimFirstMsgFromHistory(msg)
		} else if !isServiceUnavailableError(err) {
			break
		}

		sleepBeforeRetry(i)
	}

	if len(res.Choices) < 1 {
		return "", errors.New("empty response")
	}

	return res.Choices[0].Message.Content, err
}

func trimFirstMsgFromHistory(msg []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	for i, m := range msg {
		if m.Role != "ChatMessageRoleSystem" {
			return append(msg[:i], msg[i+1:]...)
		}
	}
	return msg
}
