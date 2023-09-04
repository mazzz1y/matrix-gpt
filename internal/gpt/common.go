package gpt

import (
	"time"

	"github.com/sashabaranov/go-openai"
)

func sleepBeforeRetry(i int) {
	time.Sleep(time.Duration(i*3) * time.Second)
}

func isServiceUnavailableError(err error) bool {
	e, ok := isOpenAIError(err)
	if !ok {
		return false
	}

	if e.HTTPStatusCode == 503 {
		return true
	}

	return false
}

func isTokenExceededError(err error) bool {
	e, ok := isOpenAIError(err)
	if !ok {
		return false
	}

	if s, ok := e.Code.(string); ok && s == "context_length_exceeded" {
		return true
	}

	return false
}

func isOpenAIError(err error) (*openai.APIError, bool) {
	e, ok := err.(*openai.APIError)
	return e, ok
}
