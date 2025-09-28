package conversation

import (
	"fmt"
	"sync"
)

type InMemoryConversationStore struct {
	mu    *sync.RWMutex
	table map[int]*Conversation
}

func NewDefaultConversationStore() *InMemoryConversationStore {
	return &InMemoryConversationStore{
		mu: &sync.RWMutex{},
		table: map[int]*Conversation{},
	}
}

func (cm *InMemoryConversationStore) Store(userId int, conv *Conversation) error {
	cm.mu.Lock()
	cm.table[userId] = conv
	cm.mu.Unlock()
	return nil
}

func (cm *InMemoryConversationStore) Get(userId int) (*Conversation, error) {
	cm.mu.RLock()
	conv, ok := cm.table[userId]
	cm.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no conversation found for userId %d", userId)
	}
	return conv, nil
}

func (cm *InMemoryConversationStore) Remove(userId int) error {
	cm.mu.Lock()
	delete(cm.table, userId)
	cm.mu.Unlock()
	return nil
}
