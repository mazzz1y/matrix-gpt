package gpt

import (
	"sync"
	"time"
)

// User represents a GPT user with a chat history and last message timestamp.
type User struct {
	sync.RWMutex
	History  *HistoryManager
	GPTMutex sync.Mutex
	lastMsg  time.Time
}

// NewGptUser creates a new GPT user instance with a given history size.
func NewGptUser(historySize int) *User {
	return &User{
		History: NewHistoryManager(historySize),
		lastMsg: time.Now(),
	}
}

// GetLastMsgTime retrieves the timestamp of the user's last message.
func (u *User) GetLastMsgTime() time.Time {
	u.RLock()
	defer u.RUnlock()

	return u.lastMsg
}

// UpdateLastMsgTime updates the timestamp of the user's last message to the current time.
func (u *User) UpdateLastMsgTime() {
	u.Lock()
	defer u.Unlock()

	u.lastMsg = time.Now()
}
