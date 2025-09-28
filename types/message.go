package types

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Url    string `json:"url,omitempty"`
	User   *User  `json:"user,omitempty"`
}

type Message struct {
	Id       int             `json:"message_id"`
	Date     int64           `json:"date"`
	Chat     *Chat           `json:"chat"`
	From     *User           `json:"from,omitempty"`
	Document *Document       `json:"document,omitempty"`
	Text     string          `json:"text,omitempty"`
	Entities []MessageEntity `json:"entities,omitempty"`
	ReplyTo  *Message        `json:"reply_to_message,omitempty"`
	EditDate uint            `json:"edit_date,omitempty"`
}

type MessageId struct {
	MessageId int `json:"message_id"`
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

func (m *Message) Time() time.Time {
	return time.Unix(m.Date, 0)
}

func (m *Message) IsCommand() bool {
	if len(m.Entities) > 0 {
		return m.Entities[0].IsCommand()
	}
	txt := m.Text
	if idx := strings.IndexByte(txt, '/'); idx == 0 {
		return (len(txt) > 1)
	}
	return false
}

func (m *Message) Command() (string, bool) {
	if !m.IsCommand() {
		return "", false
	}
	var command string
	if len(m.Entities) > 0 {
		e := m.Entities[0]
		command = m.Text[1:e.Length]
		if idx := strings.IndexByte(command, '@'); idx != -1 {
			command = command[:idx]
		}
		goto ret
	}
	command = m.Text[1:]
	if idx := strings.IndexByte(command, ' '); idx != -1 {
		command = command[:idx]
	}
ret:
	return command, true
}
