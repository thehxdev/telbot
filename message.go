package telbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"strings"
	"time"
)

// type MessageId struct {
// 	Id int `json:"message_id"`
// }

type MessageEntity struct {
	// Type of the entity. Can be:
	//  "mention" (@username),
	//  "hashtag" (#hashtag),
	//  "cashtag" ($USD),
	//  "bot_command" (/start@jobs_bot),
	//  "url" (https://telegram.org),
	//  "email" (do-not-reply@telegram.org),
	//  "phone_number" (+1-212-555-0123),
	//  "bold" (bold text),
	//  "italic" (italic text),
	//  "underline" (underlined text),
	//  "strikethrough" (strikethrough text),
	//  "spoiler" (spoiler message),
	//  "code" (monowidth string),
	//  "pre" (monowidth block),
	//  "text_link" (for clickable text URLs),
	//  "text_mention" (for users without usernames)
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Url    string `json:"url,omitempty"`
	User   *User  `json:"user,omitempty"`
}

type Message struct {
	Id       int             `json:"message_id"`
	Date     uint            `json:"date"`
	Chat     *Chat           `json:"chat"`
	From     *User           `json:"from,omitempty"`
	Document *Document       `json:"document,omitempty"`
	Text     string          `json:"text,omitempty"`
	Entities []MessageEntity `json:"entities,omitempty"`
	ReplyTo  *Message        `json:"reply_to_message,omitempty"`
	EditDate uint            `json:"edit_date,omitempty"`
}

type ReplyParameters struct {
	MessageId int `json:"message_id"`
	ChatId    int `json:"chat_id,omitempty"`
}

type TextMessageParams struct {
	ChatId      int              `json:"chat_id"`
	Text        string           `json:"text"`
	ReplyParams *ReplyParameters `json:"reply_parameters,omitempty"`

	// This field is not used anymore (I assume it's legacy. Use ReplyParams instead)
	ReplyToMsgId int `json:"reply_to_message_id,omitempty"`
}

func (e *MessageEntity) IsCommand() bool {
	return e.Offset == 0 && e.Type == "bot_command"
}

func (e *MessageEntity) ParseURL() (*url.URL, error) {
	if e.Url == "" {
		return nil, errors.New("Bad Url")
	}
	return url.Parse(e.Url)
}

func (tm *TextMessageParams) ToReader() (io.Reader, error) {
	b, err := json.Marshal(tm)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (tm *TextMessageParams) ContentType() string {
	return ContentTypeApplicationJson
}

func (m *Message) Time() time.Time {
	return time.Unix(int64(m.Date), 0)
}

func (m *Message) IsCommand() bool {
	if m.Entities == nil || len(m.Entities) == 0 {
		return false
	}
	e := m.Entities[0]
	return e.IsCommand()
}

func (m *Message) Command() string {
	if !m.IsCommand() {
		return ""
	}
	e := m.Entities[0]
	command := m.Text[1:e.Length]
	if i := strings.Index(command, "@"); i != 1 {
		command = command[:i]
	}
	return command
}
