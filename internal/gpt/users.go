package gpt

import (
	"time"
)

// User represents a GPT user with a chat history and last message timestamp.
type User struct {
	History *HistoryManager
	LastMsg time.Time
}

// NewGptUser creates a new GPT user instance with a given history size.
func NewGptUser(historySize int) *User {
	return &User{
		History: NewHistoryManager(historySize),
	}
}
