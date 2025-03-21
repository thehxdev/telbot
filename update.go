package telbot

import (
	"context"
	"github.com/thehxdev/telbot/types"
)

const defaultInvalidId = -1

type Update struct {
	Id      int            `json:"update_id"`
	Message *types.Message `json:"message,omitempty"`

	bot *Bot `json:"-"`
}

func (u *Update) GetBot() *Bot {
	return u.bot
}

func (u *Update) SendMessage(params TextMessageParams) (*types.Message, error) {
	return u.bot.SendMessage(params)
}

func (u *Update) UploadFile(ctx context.Context, params UploadParams, file FileInfo) (*types.Message, error) {
	return u.bot.UploadFile(ctx, params, file)
}

func (u *Update) EditMessage(params EditMessageTextParams) (*types.Message, error) {
	return u.bot.EditMessageText(params)
}

func (u *Update) DeleteMessage(chatId, messageId int) error {
	return u.bot.DeleteMessage(chatId, messageId)
}

func (u *Update) GetFile(fileId string) (*types.File, error) {
	return u.bot.GetFile(fileId)
}

func (u *Update) ChatId() int {
	if u.Message.Chat != nil {
		return u.Message.Chat.Id
	}
	return defaultInvalidId
}

func (u *Update) UserId() int {
	if u.Message.From != nil {
		return u.Message.From.Id
	}
	return defaultInvalidId
}

func (u *Update) MessageId() int {
	if u.Message != nil {
		return u.Message.Id
	}
	return defaultInvalidId
}

func (u *Update) ChatType() string {
	if u.Message.Chat != nil {
		return u.Message.Chat.Type
	}
	return ""
}
