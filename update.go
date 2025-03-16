package telbot

import (
	"bytes"
	"encoding/json"
	"io"
)

type Update struct {
	Id      int      `json:"update_id"`
	Message *Message `json:"message"`
}

type UpdateParams struct {
	Offset         int
	Limit          int
	Timeout        int
	AllowedUpdates []string
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
