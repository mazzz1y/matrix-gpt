package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"time"
)

type Gpt struct {
	client       *openai.Client
	ctx          context.Context
	model        string
	historyLimit int
	gptTimeout   time.Duration
	maxAttempts  int
	users        map[string]*User
}

// New initializes a Gpt instance with the provided configurations.
func New(token, gptModel string, historyLimit, gptTimeout, maxAttempts int, userIDs []string) *Gpt {
	users := make(map[string]*User)
	for _, id := range userIDs {
		users[id] = NewGptUser(historyLimit)
	}

	return &Gpt{
		client:       openai.NewClient(token),
		ctx:          context.Background(),
		model:        gptModel,
		historyLimit: historyLimit,
		gptTimeout:   time.Duration(gptTimeout) * time.Second,
		users:        users,
		maxAttempts:  maxAttempts,
	}
}

// GetUser retrieves the User instance associated with the given ID.
func (g *Gpt) GetUser(id string) (u *User, ok bool) {
	u, ok = g.users[id]
	return
}

// GetModel returns the GPT model string.
func (g *Gpt) GetModel() string {
	return g.model
}

// GetTimeout returns the timeout value for the GPT client.
func (g *Gpt) GetTimeout() time.Duration {
	return g.gptTimeout
}

// GetHistoryLimit returns the history limit value.
func (g *Gpt) GetHistoryLimit() int {
	return g.historyLimit
}
