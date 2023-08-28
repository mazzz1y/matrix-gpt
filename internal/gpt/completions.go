package gpt

import (
	"context"
	"errors"
	"time"

	"github.com/sashabaranov/go-openai"
)

// CreateCompletion retrieves a completion from GPT using the given user's message.
func (g *Gpt) CreateCompletion(ctx context.Context, history []openai.ChatCompletionMessage, userMsg string) ([]openai.ChatCompletionMessage, error) {
	// Append the user's message to the existing history.
	messageHistory := append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMsg,
	})

	response, err := g.createCompletionWithTimeout(ctx, messageHistory)
	if err != nil {
		return []openai.ChatCompletionMessage{}, err
	}

	messageHistory = append(messageHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: response.Choices[0].Message.Content,
	})

	return messageHistory, err
}

// createCompletionWithTimeout makes a request to get a GPT completion with a specified timeout.
func (g *Gpt) createCompletionWithTimeout(ctx context.Context, msg []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	var response openai.ChatCompletionResponse
	var err error

	ctx, cancel := context.WithTimeout(ctx, g.gptTimeout)
	defer cancel()

	for i := 0; i < g.maxAttempts; i++ {
		response, err = g.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    g.model,
				Messages: msg,
			},
		)
		if ctx.Err() == context.Canceled {
			return openai.ChatCompletionResponse{}, ctx.Err()
		}
		if err == nil {
			if len(response.Choices) < 1 || response.Choices[0].Message.Content == "" {
				err = errors.New("empty response")
			} else {
				return response, nil
			}
		}

		if isTokenExceededError(err) {
			msg = trimFirstMsgFromHistory(msg)
			continue
		}

		time.Sleep(time.Duration(i*2) * time.Second)
	}

	return openai.ChatCompletionResponse{}, err
}

func isTokenExceededError(err error) bool {
	e, ok := err.(*openai.APIError)
	if !ok {
		return false
	}

	if s, ok := e.Code.(string); ok && s == "context_length_exceeded" {
		return true
	}

	return false
}

func trimFirstMsgFromHistory(msg []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	for i, m := range msg {
		if m.Role != "ChatMessageRoleSystem" {
			return append(msg[:i], msg[i+1:]...)
		}
	}
	return msg
}
