package bot

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

// newHistoryManager initializes a HistoryManager instance with the provided size.
func newHistoryManager(size int) *historyManager {
	return &historyManager{
		Storage: make([]openai.ChatCompletionMessage, 0),
		Size:    size,
	}
}

// resetHistory clears the current chat history.
func (m *historyManager) resetHistory() {
	m.Lock()
	defer m.Unlock()

	if len(m.Storage) > 0 {
		m.Storage = make([]openai.ChatCompletionMessage, 0)
	}
}

// addMessage appends a new message to the chat history.
func (m *historyManager) updateHistory(h []openai.ChatCompletionMessage) {
	m.Lock()
	defer m.Unlock()

	m.Storage = h
	m.trimHistory()
}

// getHistory retrieves the current chat history.
func (m *historyManager) getHistory() []openai.ChatCompletionMessage {
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
