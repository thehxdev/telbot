package telbot

import "time"

const (
	defaultInvalidId        = -1
	defaultOperationTimeout = time.Second * 5
	getUpdatesSleepTime     = time.Second * 1
)

const (
	ChatTypePrivate    string = "private"
	ChatTypeGroup      string = "group"
	ChatTypeSuperGroup string = "supergroup"
	ChatTypeChannel    string = "channel"
)

const (
	ContentTypeFormUrlEncoded    = "application/x-www-form-urlencoded"
	ContentTypeMultipartFormData = "multipart/form-data"
	ContentTypeApplicationJson   = "application/json"
)

const (
	MethodGetMe           = "getMe"
	MethodGetUpdates      = "getUpdates"
	MethodSendMessage     = "sendMessage"
	MethodGetFile         = "getFile"
	MethodEditMessageText = "editMessageText"
	MethodDeleteMessage   = "deleteMessage"
)

const (
	MessageEntityTypeMention       = "mention"
	MessageEntityTypeHashtag       = "hashtag"
	MessageEntityTypeCashtag       = "cachtag"
	MessageEntityTypeBotCommand    = "bot_command"
	MessageEntityTypeUrl           = "url"
	MessageEntityTypeEmail         = "email"
	MessageEntityTypePhoneNumber   = "phone_number"
	MessageEntityTypeBold          = "bold"
	MessageEntityTypeItalic        = "italic"
	MessageEntityTypeUnderline     = "underline"
	MessageEntityTypeStrikeThrough = "strikethrough"
	MessageEntityTypeSpoiler       = "spoiler"
	MessageEntityTypeCode          = "code"
	MessageEntityTypePre           = "pre"
	MessageEntityTypeTextLink      = "text_link"
	MessageEntityTypeTextMention   = "text_mention"
)
