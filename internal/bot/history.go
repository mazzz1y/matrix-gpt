package bot

import (
	"sync"

	"github.com/sashabaranov/go-openai"
)

// historyManager manages chat histories for GPT interactions.
type historyManager struct {
	sync.RWMutex
	storage []openai.ChatCompletionMessage
	maxSize int
}

// newHistoryManager initializes a HistoryManager instance with the provided size.
func newHistoryManager(maxSize int) *historyManager {
	return &historyManager{
		storage: make([]openai.ChatCompletionMessage, 0),
		maxSize: maxSize,
	}
}

// reset clears the current chat history.
func (m *historyManager) reset() {
	m.Lock()
	defer m.Unlock()

	if len(m.storage) > 0 {
		m.storage = make([]openai.ChatCompletionMessage, 0)
	}
}

// save keeps the last 'm.Size' messages in memory.
func (m *historyManager) save(h []openai.ChatCompletionMessage) {
	m.Lock()
	defer m.Unlock()

	if m.maxSize != 0 && len(h) > m.maxSize {
		m.storage = h[len(h)-m.maxSize:]
	} else {
		m.storage = h
	}
}

// get retrieves the current chat history.
func (m *historyManager) get() []openai.ChatCompletionMessage {
	m.RLock()
	defer m.RUnlock()

	return m.storage
}

// getSize retrieves the current history size.
func (m *historyManager) getSize() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.storage)
}
