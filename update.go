package telbot

import (
	"context"
	"bytes"
	"encoding/json"
	"io"
)

type Update struct {
	Id      int      `json:"update_id"`
	Message *Message `json:"message,omitempty"`

	bot     *Bot     `json:"-"`
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

func (u *Update) GetBot() *Bot {
	return u.bot
}

func (u *Update) SendMessage(params TextMessageParams) (*Message, error) {
	return u.bot.SendMessage(params)
}

func (u *Update) UploadFile(ctx context.Context, params UploadParams, file FileInfo) (*Message, error) {
	return u.bot.UploadFile(ctx, params, file)
}
