package gpt

import (
	"github.com/sashabaranov/go-openai"
	"sync"
)

// HistoryManager manages chat histories for GPT interactions.
type HistoryManager struct {
	sync.RWMutex
	Storage []openai.ChatCompletionMessage
	Size    int
}

// NewHistoryManager initializes a HistoryManager instance with the provided size.
func NewHistoryManager(size int) *HistoryManager {
	return &HistoryManager{
		Storage: make([]openai.ChatCompletionMessage, 0),
		Size:    size,
	}
}

// ResetHistory clears the current chat history.
// Returns true if history was cleared, false otherwise.
func (m *HistoryManager) ResetHistory() bool {
	m.Lock()
	defer m.Unlock()

	if len(m.Storage) > 0 {
		m.Storage = make([]openai.ChatCompletionMessage, 0)
		return true
	}

	return false
}

// AddMessage appends a new message to the chat history.
func (m *HistoryManager) AddMessage(msgType, msgContent string) {
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
func (m *HistoryManager) GetHistory() []openai.ChatCompletionMessage {
	m.RLock()
	defer m.RUnlock()

	return m.Storage
}

// trimHistory ensures the chat history doesn't exceed its size limit.
func (m *HistoryManager) trimHistory() {
	if len(m.Storage) > m.Size {
		m.Storage = m.Storage[len(m.Storage)-m.Size:]
	}
}
