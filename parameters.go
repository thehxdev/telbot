package telbot

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	"github.com/thehxdev/telbot/types"
)

type IFileInfo interface {
	UploadInfo() (string, io.Reader, error)
	FileKind() string
}

type FileReader struct {
	io.Reader
	Kind     string
	FileName string
}

func (fr *FileReader) UploadInfo() (string, io.Reader, error) {
	return fr.FileName, fr.Reader, nil
}

func (fr *FileReader) FileKind() string {
	return fr.Kind
}

type UpdateParams struct {
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type ReplyParameters struct {
	MessageId int `json:"message_id"`
	ChatId    int `json:"chat_id,omitempty"`
}

type SuggestedPostParameters struct {
	Price    *types.SuggestedPostPrice `json:"price,omitempty"`
	SendDate int                       `json:"send_date,omitempty"`
}

type IReplyMarkup any

type TextMessageParams struct {
	BusinessConnectionId    string                    `json:"business_connection_id,omitempty"`
	ChatId                  int                       `json:"chat_id"`
	Text                    string                    `json:"text"`
	ParseMode               string                    `json:"parse_mode,omitempty"`
	Entities                []types.MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions      *types.LinkPreviewOptions `json:"link_preview_options,omitempty"`
	DisableNotification     bool                      `json:"disable_notification,omitempty"`
	ProtectContent          bool                      `json:"protect_content,omitempty"`
	AllowPaidBroadcast      bool                      `json:"allow_paid_broadcast,omitempty"`
	MessageEffectId         string                    `json:"message_effect_id,omitempty"`
	SuggestedPostParameters *SuggestedPostParameters  `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters          `json:"reply_parameters,omitempty"`

	// Must be one of "InlineKeyboardMarkup", "ReplyKeyboardMarkup", "ReplyKeyboardRemove"
	// or "ForceReply" types
	ReplyMarkup IReplyMarkup `json:"reply_markup,omitempty"`

	// This field is not used anymore (I assume it's legacy. Use ReplyParams instead)
	// Included for backward compatibility.
	ReplyToMessageId int `json:"reply_to_message_id,omitempty"`
}

type EditMessageTextParams struct {
	ChatId    int                   `json:"chat_id"`
	MessageId int                   `json:"message_id"`
	Text      string                `json:"text"`
	ParseMode string                `json:"parse_mode,omitempty"`
	Entities  []types.MessageEntity `json:"entities,omitempty"`
}

type UploadParams struct {
	ChatId int    `json:"chat_id"`
	Method string `json:"-"`
}

func (up *UploadParams) ToStringMap() (map[string]string, error) {
	p := map[string]string{}
	if up.ChatId != 0 {
		p["chat_id"] = strconv.Itoa(up.ChatId)
	}
	return p, nil
}

func ParamsToReader(params any) (io.Reader, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
