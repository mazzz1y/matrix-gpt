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

// reset clears the current chat history.
func (m *historyManager) reset() {
	m.Lock()
	defer m.Unlock()

	if len(m.Storage) > 0 {
		m.Storage = make([]openai.ChatCompletionMessage, 0)
	}
}

// save keeps the last 'm.Size' messages in memory.
func (m *historyManager) save(h []openai.ChatCompletionMessage) {
	m.Lock()
	defer m.Unlock()

	if len(h) > m.Size {
		m.Storage = h[len(h)-m.Size:]
	} else {
		m.Storage = h
	}
}

// get retrieves the current chat history.
func (m *historyManager) get() []openai.ChatCompletionMessage {
	m.RLock()
	defer m.RUnlock()

	return m.Storage
}
