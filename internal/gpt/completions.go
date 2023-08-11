package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"time"
)

// CreateCompletion retrieves a completion from GPT using the given user's message.
func (g *Gpt) CreateCompletion(u *User, userMsg string) (string, error) {
	// Append the user's message to the existing history.
	messageHistory := append(u.History.GetHistory(), openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMsg,
	})

	var (
		response openai.ChatCompletionResponse
		err      error
	)

	// Try creating a completion up to the maximum number of allowed attempts.
	for i := 0; i < g.maxAttempts; i++ {
		response, err = g.createCompletionWithTimeout(messageHistory)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	// Update the user's history with both the user message and the assistant's response.
	u.History.AddMessage(openai.ChatMessageRoleUser, userMsg)
	u.History.AddMessage(openai.ChatMessageRoleAssistant, response.Choices[0].Message.Content)

	return response.Choices[0].Message.Content, err
}

// createCompletionWithTimeout makes a request to get a GPT completion with a specified timeout.
func (g *Gpt) createCompletionWithTimeout(msg []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	ctx, cancel := context.WithTimeout(g.ctx, g.gptTimeout)
	defer cancel()

	return g.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    g.model,
			Messages: msg,
		},
	)
}
