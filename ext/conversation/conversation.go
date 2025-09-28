package conversation

import (
	"fmt"
	"sync"

	"github.com/thehxdev/telbot"
)

type Conversation struct {
	Next   ConversationHandler
	ChatId int
	UserId int
}

type ConversationHandler func(*Conversation, telbot.Update) error

type EndConversation struct{}

func (e *EndConversation) Error() string {
	return "end conversation"
}

var convMap = &struct {
	mu    *sync.RWMutex
	table map[int]*Conversation
}{
	mu:    &sync.RWMutex{},
	table: make(map[int]*Conversation),
}

func HasConversation(chatId, userId int) bool {
	convMap.mu.RLock()
	defer convMap.mu.RUnlock()
	if conv, ok := convMap.table[userId]; ok {
		return conv.ChatId == chatId
	}
	return false
}

func CallNext(update telbot.Update) error {
	userId := update.Message.From.Id

	convMap.mu.RLock()
	conv, ok := convMap.table[userId]
	convMap.mu.RUnlock()

	if !ok {
		return fmt.Errorf("no conversation found for userId %d", userId)
	}

	err := conv.Next(conv, update)
	switch err.(type) {
	case *EndConversation:
		convMap.mu.Lock()
		delete(convMap.table, conv.UserId)
		convMap.mu.Unlock()
		err = nil
	}

	return err
}

func Start(startHandler ConversationHandler, update telbot.Update) {
	userId := update.Message.From.Id
	c := &Conversation{
		UserId: userId,
		ChatId: update.Message.Chat.Id,
	}

	convMap.mu.Lock()
	convMap.table[userId] = c
	convMap.mu.Unlock()

	startHandler(c, update)
}
