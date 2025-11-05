package telbot

import (
	"github.com/thehxdev/telbot/types"
)

// TODO: Complete the `Update` type
type Update struct {
	Id                     int                                `json:"update_id"`
	Message                *types.Message                     `json:"message,omitempty"`
	EditedMessage          *types.Message                     `json:"edited_message,omitempty"`
	ChannelPost            *types.Message                     `json:"channel_post,omitempty"`
	EditedChannelPost      *types.Message                     `json:"edited_cahannel_post,omitempty"`
	BusinessConnection     *types.BusinessConnection          `json:"business_connection,omitempty"`
	BusinessMessage        *types.Message                     `json:"business_message,omitempty"`
	EditedBusinessMessage  *types.Message                     `json:"edited_business_message,omitempty"`
	DeletedBusinessMessage *types.BusinessMessageDeleted      `json:"deleted_business_messages,omitempty"`
	MessageReaction        *types.MessageReactionUpdated      `json:"message_reaction,omitempty"`
	MessageReactionCount   *types.MessageReactionCountUpdated `json:"message_reaction_count,omitempty"`
	InlineQuery            *types.InlineQuery                 `json:"inline_query,omitempty"`
	ChosenInlineResult     *types.ChosenInlineResult          `json:"chosen_inline_result,omitempty"`
	CallbackQuery          *types.CallbackQuery               `json:"callback_query,omitempty"`
	ShippingQuery          *types.ShippingQuery               `json:"shipping_query,omitempty"`
	PreCheckoutQuery       *types.PreCheckoutQuery            `json:"pre_checkout_query,omitempty"`
	PurchasedPaidMedia     *types.PaidMediaPurchased          `json:"purchased_paid_media,omitempty"`

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
