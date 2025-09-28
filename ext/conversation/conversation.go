package conversation

import (
	"github.com/thehxdev/telbot"
)

type Conversation struct {
	Next   ConversationHandler
	ChatId int
	UserId int
}

type ConversationStore interface {
	Store(userId int, conv *Conversation) error
	Get(userId int) (*Conversation, error)
	Remove(userId int) error
}

type ConversationHandler func(*Conversation, telbot.Update) error

type EndConversation struct{}

func (e *EndConversation) Error() string {
	return "end conversation"
}

var convStore ConversationStore = NewDefaultConversationStore()

func SetConversationStore(cs ConversationStore) {
	convStore = cs
}

func HasConversation(chatId, userId int) bool {
	if conv, err := convStore.Get(userId); err == nil {
		return conv.ChatId == chatId
	}
	return false
}

func CallNext(update telbot.Update) error {
	userId := update.Message.From.Id

	conv, err := convStore.Get(userId)
	if err != nil {
		return err
	}

	err = conv.Next(conv, update)
	switch err.(type) {
	case *EndConversation:
		convStore.Remove(userId)
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
	_ = convStore.Store(userId, c)
	startHandler(c, update)
}
