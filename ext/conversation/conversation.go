package conversation

import (
	"errors"
	"sync"

	"github.com/thehxdev/telbot"
)

type ErrRepeatStage struct{}

type ErrEnd struct{}

var (
	EndConversation = &ErrEnd{}
	RepeatSameStage = &ErrRepeatStage{}
)

func (e *ErrRepeatStage) Error() string {
	return "repeat same stage"
}

func (e *ErrEnd) Error() string {
	return "end conversation"
}

type Conversation struct {
	stages  []telbot.UpdateHandlerFunc
	count   int
	current int
	chatId  int
	userId  int
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
		return conv.chatId == chatId
	}
	return false
}

func HandleUpdate(update telbot.Update) error {
	userId := update.Message.From.Id

	convMap.mu.RLock()
	conv := convMap.table[userId]
	convMap.mu.RUnlock()

	err := conv.callStage(update)
	if err != nil {
		switch err.(type) {
		case *ErrEnd:
			endConversation(userId)
		case *ErrRepeatStage:
			return nil
		default:
			return err
		}
	}

	conv.current += 1
	return nil
}

func New(stages []telbot.UpdateHandlerFunc) (Conversation, error) {
	stagesLen := len(stages)
	if stagesLen == 0 {
		return Conversation{}, errors.New("empty stage handlers list")
	}

	conv := Conversation{
		stages:  []telbot.UpdateHandlerFunc{},
		count:   stagesLen,
		current: 0,
	}
	conv.stages = append(conv.stages, stages...)

	return conv, nil
}

func (conv Conversation) Start(update telbot.Update) {
	userId := update.Message.From.Id
	c := new(Conversation)
	*c = conv

	c.userId = userId
	c.chatId = update.Message.Chat.Id

	convMap.mu.Lock()
	convMap.table[userId] = c
	convMap.mu.Unlock()

	HandleUpdate(update)
}

func (conv *Conversation) getStageHandler() (telbot.UpdateHandlerFunc, error) {
	if conv.current < conv.count {
		return conv.stages[conv.current], nil
	}
	return nil, &ErrEnd{}
}

func (conv *Conversation) callStage(update telbot.Update) error {
	fn, err := conv.getStageHandler()
	if err != nil {
		return err
	}
	return fn(update)
}

func endConversation(userId int) {
	convMap.mu.Lock()
	delete(convMap.table, userId)
	convMap.mu.Unlock()
}
