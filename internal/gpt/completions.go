package gpt

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// CreateCompletion retrieves a completion from GPT using the given user's message.
func (g *Gpt) CreateCompletion(ctx context.Context, u *User, userMsg string) (string, error) {
	// Append the user's message to the existing history.
	messageHistory := append(u.History.GetHistory(), openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMsg,
	})

	var response openai.ChatCompletionResponse
	var err error
	response, err = g.createCompletionWithTimeout(ctx, messageHistory)
	if err != nil {
		return "", err
	}

	// Update the user's history with both the user message and the assistant's response.
	u.History.AddMessage(openai.ChatMessageRoleUser, userMsg)
	u.History.AddMessage(openai.ChatMessageRoleAssistant, response.Choices[0].Message.Content)

	return response.Choices[0].Message.Content, err
}

// createCompletionWithTimeout makes a request to get a GPT completion with a specified timeout.
func (g *Gpt) createCompletionWithTimeout(ctx context.Context, msg []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	var err error
	for i := 0; i < g.maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(ctx, g.gptTimeout)
		defer cancel()

		response, err := g.client.CreateChatCompletion(
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
			return response, nil
		}
	}

	return openai.ChatCompletionResponse{}, err
}
