package telbot

import (
	"github.com/thehxdev/telbot/types"
)

const defaultInvalidId = -1

type Update struct {
	Id                int            `json:"update_id"`
	Message           *types.Message `json:"message,omitempty"`
	EditedMessage     *types.Message `json:"edited_message,omitempty"`
	ChannelPost       *types.Message `json:"channel_post,omitempty"`
	EditedChannelPost *types.Message `json:"edited_cahannel_post,omitempty"`

	Bot *Bot `json:"-"`
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
