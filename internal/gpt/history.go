package gpt

import (
	"sync"

	"github.com/sashabaranov/go-openai"
)

// historyManager manages chat histories for GPT interactions.
type historyManager struct {
	sync.RWMutex
	Storage []openai.ChatCompletionMessage
	Size    int
}

// NewHistoryManager initializes a HistoryManager instance with the provided size.
func newHistoryManager(size int) *historyManager {
	return &historyManager{
		Storage: make([]openai.ChatCompletionMessage, 0),
		Size:    size,
	}
}

// ResetHistory clears the current chat history.
func (m *historyManager) ResetHistory() {
	m.Lock()
	defer m.Unlock()

	if len(m.Storage) > 0 {
		m.Storage = make([]openai.ChatCompletionMessage, 0)
	}
}

// AddMessage appends a new message to the chat history.
func (m *historyManager) AddMessage(msgType, msgContent string) {
	m.Lock()
	defer m.Unlock()

	message := openai.ChatCompletionMessage{
		Role:    msgType,
		Content: msgContent,
	}
	m.Storage = append(m.Storage, message)
	m.trimHistory()
}

// GetHistory retrieves the current chat history.
func (m *historyManager) GetHistory() []openai.ChatCompletionMessage {
	m.RLock()
	defer m.RUnlock()

	return m.Storage
}

// trimHistory ensures the chat history doesn't exceed its size limit.
func (m *historyManager) trimHistory() {
	if len(m.Storage) > m.Size {
		m.Storage = m.Storage[len(m.Storage)-m.Size:]
	}
}
