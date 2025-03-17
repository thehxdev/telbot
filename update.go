package telbot

import (
	"bytes"
	"encoding/json"
	"io"
)

type Update struct {
	Id      int      `json:"update_id"`
	Message *Message `json:"message,omitempty"`
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
