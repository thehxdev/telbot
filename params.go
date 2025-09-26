package telbot

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	"github.com/thehxdev/telbot/types"
)

type FileInfo interface {
	UploadInfo() (string, io.Reader, error)
	FileKind() string
}

type FileReader struct {
	Kind     string
	FileName string
	Reader   io.Reader
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

func (up *UpdateParams) ToReader() (io.Reader, error) {
	b, err := json.Marshal(up)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (up *UpdateParams) ContentType() string {
	return ContentTypeApplicationJson
}

type ReplyParameters struct {
	MessageId int `json:"message_id"`
	ChatId    int `json:"chat_id,omitempty"`
}

type TextMessageParams struct {
	ChatId      int                   `json:"chat_id"`
	Text        string                `json:"text"`
	ParseMode   string                `json:"parse_mode,omitempty"`
	ReplyParams *ReplyParameters      `json:"reply_parameters,omitempty"`
	Entities    []types.MessageEntity `json:"entities,omitempty"`

	// This field is not used anymore (I assume it's legacy. Use ReplyParams instead)
	// Included for backward compatibility.
	ReplyToMsgId int `json:"reply_to_message_id,omitempty"`
}

func (p *TextMessageParams) ToReader() (io.Reader, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (p *TextMessageParams) ContentType() string {
	return ContentTypeApplicationJson
}

type EditMessageTextParams struct {
	ChatId    int                   `json:"chat_id"`
	MessageId int                   `json:"message_id"`
	Text      string                `json:"text"`
	ParseMode string                `json:"parse_mode,omitempty"`
	Entities  []types.MessageEntity `json:"entities,omitempty"`
}

func (p *EditMessageTextParams) ToReader() (io.Reader, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (p *EditMessageTextParams) ContentType() string {
	return ContentTypeApplicationJson
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
