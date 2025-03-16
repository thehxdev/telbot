package telbot

import (
	"io"
	"strconv"
)

type File struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

type Document struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

type UploadParams struct {
	ChatId int    `json:"chat_id"`
	Method string `json:"-"`
}

type FileInfo interface {
	// NeedUpload() bool
	UploadInfo() (string, io.Reader, error)
	FileKind() string
}

type FileReader struct {
	Kind     string
	FileName string
	Reader   io.Reader
}

func (up *UploadParams) ToStringMap() (StringMap, error) {
	p := StringMap{}
	if up.ChatId != 0 {
		p["chat_id"] = strconv.Itoa(up.ChatId)
	}
	return p, nil
}

func (fr FileReader) UploadInfo() (string, io.Reader, error) {
	return fr.FileName, fr.Reader, nil
}

func (fr FileReader) FileKind() string {
	return fr.Kind
}
