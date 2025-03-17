// A simple telegram bot api library that supports a little subset of official telegram api
// for building small and simple bots.
// Inspired by `https://github.com/go-telegram-bot-api/telegram-bot-api`

// FIXME: Don't use context.Background() for every request

package telbot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

const GetUpdatesSleepTime = time.Second * 1

const (
	ContentTypeFormUrlEncoded    = "application/x-www-form-urlencoded"
	ContentTypeMultipartFormData = "multipart/form-data"
	ContentTypeApplicationJson   = "application/json"
)

const (
	MethodGetMe       = "getMe"
	MethodGetUpdates  = "getUpdates"
	MethodSendMessage = "sendMessage"
	MethodGetFile     = "getFile"
)

type Bot struct {
	Token       string
	BaseUrl     string
	BaseFileUrl string
	Self        *User
	UpdatesChan chan Update
	Client      *http.Client

	// Shutdown signal channel
	sdChan chan struct{}
}

type UpdateHandlerFunc func(bot *Bot, update Update) error

type StringMap map[string]string

type RequestBody interface {
	ToReader() (io.Reader, error)
	ContentType() string
}

type WrappedReader struct {
	io.Reader
	// Content Type
	CType string
}

func (wr *WrappedReader) ToReader() (io.Reader, error) {
	return wr.Reader, nil
}

func (wr *WrappedReader) ContentType() string {
	return wr.CType
}

type RequestInfo struct {
	Method      string
	BaseUrl     string
	HasBody     bool
	Body        RequestBody
	ContentType string
}

type APIResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
}

// Create a new instance of Bot
func New(token string, host ...string) (*Bot, error) {
	h := "api.telegram.org"
	if len(host) > 0 {
		h = host[0]
	}

	b := &Bot{
		Token:       token,
		BaseUrl:     fmt.Sprintf("https://%s/bot%s", h, token),
		BaseFileUrl: fmt.Sprintf("https://%s/file/bot%s", h, token),
		Client:      &http.Client{},
		sdChan:      make(chan struct{}),
	}

	botUser, err := b.GetMe()
	if err != nil {
		return nil, err
	}
	b.Self = botUser

	return b, nil
}

func (b *Bot) Shutdown() {
	close(b.sdChan)
	close(b.UpdatesChan)
}

func CreateMethodUrl(baseUrl string, method string) string {
	return fmt.Sprintf("%s/%s", baseUrl, method)
}

func (b *Bot) SendRequest(ctx context.Context, info RequestInfo) (*APIResponse, error) {
	var err error
	var reqBody io.Reader = nil
	if info.HasBody {
		reqBody, err = info.Body.ToReader()
		if err != nil {
			return nil, err
		}
	}

	reqUrl := CreateMethodUrl(info.BaseUrl, info.Method)
	req, err := http.NewRequestWithContext(ctx, "POST", reqUrl, reqBody)
	if err != nil {
		return nil, err
	}
	if info.ContentType != "" {
		req.Header.Add("Content-Type", info.ContentType)
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	apiResp := &APIResponse{}
	if err := json.NewDecoder(resp.Body).Decode(apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Ok {
		return nil, errors.New(apiResp.Description)
	}

	return apiResp, nil
}

func (b *Bot) GetMe() (*User, error) {
	u := &User{}
	apiResp, err := b.SendRequest(context.Background(), RequestInfo{
		Method:      MethodGetMe,
		BaseUrl:     b.BaseUrl,
		HasBody:     false,
		Body:        nil,
		ContentType: "",
	})
	if err != nil {
		return nil, err
	}

	if !apiResp.Ok {
		return nil, errors.New(apiResp.Description)
	}

	err = json.Unmarshal(apiResp.Result, u)
	return u, err
}

func (b *Bot) GetUpdates(params UpdateParams) ([]Update, error) {
	updates := []Update{}
	apiResp, err := b.SendRequest(context.Background(), RequestInfo{
		Method:      MethodGetUpdates,
		BaseUrl:     b.BaseUrl,
		HasBody:     true,
		Body:        &params,
		ContentType: params.ContentType(),
	})
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(apiResp.Result, &updates)
	if err != nil {
		return nil, err
	}
	return updates, nil
}

func (b *Bot) UploadFile(ctx context.Context, params UploadParams, file FileInfo) (*Message, error) {
	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)

	go func() {
		defer pipeWriter.Close()
		defer multipartWriter.Close()

		pMap, _ := params.ToStringMap()
		for key, value := range pMap {
			if err := multipartWriter.WriteField(key, value); err != nil {
				pipeWriter.CloseWithError(err)
				return
			}
		}

		fileName, fileReader, err := file.UploadInfo()
		if err != nil {
			pipeWriter.CloseWithError(err)
			return
		}
		part, err := multipartWriter.CreateFormFile(file.FileKind(), fileName)
		if err != nil {
			pipeWriter.CloseWithError(err)
			return
		}

		if _, err := io.Copy(part, fileReader); err != nil {
			pipeWriter.CloseWithError(err)
			return
		}

		if closer, ok := fileReader.(io.ReadCloser); ok {
			if err = closer.Close(); err != nil {
				pipeWriter.CloseWithError(err)
				return
			}
		}
	}()

	apiResp, err := b.SendRequest(context.Background(), RequestInfo{
		Method:      "sendDocument",
		BaseUrl:     b.BaseUrl,
		HasBody:     true,
		Body:        &WrappedReader{Reader: pipeReader},
		ContentType: multipartWriter.FormDataContentType(),
	})
	if err != nil {
		return nil, err
	}

	msg := &Message{}
	if err := json.NewDecoder(bytes.NewReader(apiResp.Result)).Decode(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (b *Bot) GetFile(fileId string) (*File, error) {
	jbytes, err := json.Marshal(StringMap{"file_id": fileId})
	apiResp, err := b.SendRequest(context.Background(), RequestInfo{
		Method:      MethodGetFile,
		BaseUrl:     b.BaseUrl,
		HasBody:     true,
		Body:        &WrappedReader{Reader: bytes.NewReader(jbytes)},
		ContentType: ContentTypeApplicationJson,
	})
	if err != nil {
		return nil, err
	}

	file := &File{}
	err = json.Unmarshal(apiResp.Result, file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (b *Bot) SendMessage(params TextMessageParams) (*Message, error) {
	apiResp, err := b.SendRequest(context.Background(), RequestInfo{
		Method:      MethodSendMessage,
		BaseUrl:     b.BaseUrl,
		HasBody:     true,
		Body:        &params,
		ContentType: params.ContentType(),
	})
	if err != nil {
		return nil, err
	}

	msg := &Message{}
	err = json.Unmarshal(apiResp.Result, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (b *Bot) StartPolling(params UpdateParams) (<-chan Update, error) {
	b.UpdatesChan = make(chan Update, params.Limit)

	go func() {
		for {
			select {
			case <-b.sdChan:
				return
			default:
			}
			updates, err := b.GetUpdates(params)
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second * 5)
				continue
			}
			for _, update := range updates {
				if update.Id >= params.Offset {
					params.Offset = update.Id + 1
					if update.Message != nil {
						if hasConversation(update) {
							go handleConversationUpdate(update)
							continue
						}
						b.UpdatesChan <- update
					}
				}
			}
			time.Sleep(GetUpdatesSleepTime)
		}
	}()

	return b.UpdatesChan, nil
}
