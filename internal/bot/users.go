package bot

import (
	"context"
	"sync"
	"time"
)

// user represents a GPT user with a chat history and last message timestamp.
type user struct {
	sync.RWMutex
	history   *historyManager
	reqMutex  sync.Mutex
	activeReq *request
	lastMsg   time.Time
}

// request represents a real-time request from a user.
// each request is given a unique ID for tracking and a cancel function to stop the request if needed.
type request struct {
	ID     string
	cancel context.CancelFunc
}

// getLastMsgTime retrieves the timestamp of the user's last message.
func (u *user) getLastMsgTime() time.Time {
	u.RLock()
	defer u.RUnlock()

	return u.lastMsg
}

// updateLastMsgTime updates the timestamp of the user's last message to the current time.
func (u *user) updateLastMsgTime() {
	u.Lock()
	defer u.Unlock()

	u.lastMsg = time.Now()
}

// createRequestContext creates a new context for a request and stores it as the active request.
func (u *user) createRequestContext(id string) *context.Context {
	u.Lock()
	defer u.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	u.activeReq = &request{
		ID:     id,
		cancel: cancel,
	}

	return &ctx
}

// getActiveRequest gets the current active request ID.
func (u *user) getActiveRequestID() (id string, exists bool) {
	u.Lock()
	defer u.Unlock()

	if u.activeReq != nil {
		return u.activeReq.ID, true
	}

	return "", false
}

// cancelRequestContext cancels the context of the request.
func (u *user) cancelRequestContext(id string) {
	u.Lock()
	defer u.Unlock()

	if u.activeReq != nil && u.activeReq.ID == id {
		u.activeReq.cancel()
		u.activeReq = nil
	}
}

// newGptUser creates a new GPT user instance with a given history size.
func newGptUser(historySize int) *user {
	return &user{
		history: newHistoryManager(historySize),
		lastMsg: time.Now(),
	}
}
