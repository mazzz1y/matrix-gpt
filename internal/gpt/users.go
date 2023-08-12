package gpt

import (
	"context"
	"sync"
	"time"
)

// User represents a GPT user with a chat history and last message timestamp.
type User struct {
	sync.RWMutex
	History   *historyManager
	ReqMutex  sync.Mutex
	activeReq *request
	lastMsg   time.Time
}

type request struct {
	ID     string
	cancel context.CancelFunc
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

// CreateRequestContext creates a new context for a request and stores it as the active request.
func (u *User) CreateRequestContext(id string) *context.Context {
	u.Lock()
	defer u.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	u.activeReq = &request{
		ID:     id,
		cancel: cancel,
	}

	return &ctx
}

// GetActiveRequest gets the current active request ID.
func (u *User) GetActiveRequestID() (id string, exists bool) {
	u.Lock()
	defer u.Unlock()

	if u.activeReq != nil {
		return u.activeReq.ID, true
	}

	return "", false
}

// CancelRequestContext cancels the context of the request.
func (u *User) CancelRequestContext(id string) {
	u.Lock()
	defer u.Unlock()

	if u.activeReq != nil && u.activeReq.ID == id {
		u.activeReq.cancel()
		u.activeReq = nil
	}
}

// newGptUser creates a new GPT user instance with a given history size.
func newGptUser(historySize int) *User {
	return &User{
		History: newHistoryManager(historySize),
		lastMsg: time.Now(),
	}
}
