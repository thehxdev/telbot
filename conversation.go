package telbot

import (
	"errors"
	"sync"
)

type ErrConversationRepeat struct{}
func (e *ErrConversationRepeat) Error() string {
	return "repeat same stage"
}

type ErrConversationEnd struct{}
func (e *ErrConversationEnd) Error() string {
	return "end conversation"
}

type Conversation struct {
	bot     *Bot
	stages  []UpdateHandlerFunc
	count   int
	current int
	chatId  int
	userId  int
}

type ConversationConfig struct {
	Bot *Bot
	Stages []UpdateHandlerFunc
}

var convMap = &struct {
	mu    *sync.RWMutex
	table map[int]*Conversation
}{
	mu:    &sync.RWMutex{},
	table: make(map[int]*Conversation),
}

func hasConversation(update Update) bool {
	if update.Message.Chat == nil || update.Message.From == nil {
		return false
	}
	userId := update.Message.From.Id
	chatId := update.Message.Chat.Id
	convMap.mu.RLock()
	defer convMap.mu.RUnlock()
	if conv, ok := convMap.table[userId]; ok {
		return conv.chatId == chatId
	}
	return false
}

func handleConversationUpdate(update Update) {
	userId := update.Message.From.Id

	convMap.mu.RLock()
	conv := convMap.table[userId]
	convMap.mu.RUnlock()

	err := conv.callStage(update)
	if err != nil {
		switch err.(type) {
		case *ErrConversationEnd:
			endConversation(userId)
		case *ErrConversationRepeat:
			return
		default:
			return
		}
	}
	conv.current += 1
}

func NewConversation(conf ConversationConfig) (Conversation, error) {
	stagesLen := len(conf.Stages)
	if stagesLen == 0 {
		return Conversation{}, errors.New("empty stage handlers list")
	}

	conv := Conversation{
		bot: conf.Bot,
		stages: []UpdateHandlerFunc{},
		count: stagesLen,
		current: 0,
	}
	conv.stages = append(conv.stages, conf.Stages...)

	return conv, nil
}

func (conv Conversation) Start(update Update) {
	userId := update.Message.From.Id
	c := new(Conversation)
	*c = conv

	c.userId = userId
	c.chatId = update.Message.Chat.Id

	convMap.mu.Lock()
	convMap.table[userId] = c
	convMap.mu.Unlock()

	handleConversationUpdate(update)
}

func (conv *Conversation) getStageHandler() (UpdateHandlerFunc, error) {
	if conv.current < conv.count {
		return conv.stages[conv.current], nil
	}
	return nil, &ErrConversationEnd{}
}

func (conv *Conversation) callStage(update Update) error {
	fn, err := conv.getStageHandler()
	if err != nil {
		return err
	}
	return fn(conv.bot, update)
}

func endConversation(userId int) {
	convMap.mu.Lock()
	delete(convMap.table, userId)
	convMap.mu.Unlock()
}
