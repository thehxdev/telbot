package telbot

// TODO: find a better way to use context package with this library.

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

	"github.com/thehxdev/telbot/types"
)

const GetUpdatesSleepTime = time.Second * 1

type Bot struct {
	Token       string
	BaseUrl     string
	BaseFileUrl string
	Self        *types.User
	Client      *http.Client

	// Shutdown signal channel
	sdChan      chan struct{}
	updatesChan chan Update
}

type UpdateHandlerFunc func(update Update) error

// type StringMap map[string]string

type RequestBody interface {
	ToReader() (io.Reader, error)
	ContentType() string
}

// This type wraps an `io.Reader` and implements `RequestBody` interface for it.
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

// Create a new Bot
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
	close(b.updatesChan)
}

func CreateMethodUrl(baseUrl string, method string) string {
	return fmt.Sprintf("%s/%s", baseUrl, method)
}

func (b *Bot) SendRequest(ctx context.Context, baseUrl string, info RequestInfo) (*APIResponse, error) {
	var err error
	var reqBody io.Reader = nil
	if info.HasBody {
		reqBody, err = info.Body.ToReader()
		if err != nil {
			return nil, err
		}
	}

	reqUrl := CreateMethodUrl(baseUrl, info.Method)
	req, err := http.NewRequestWithContext(ctx, "POST", reqUrl, reqBody)
	if err != nil {
		return nil, err
	}
	if info.ContentType != "" {
		req.Header.Set("Content-Type", info.ContentType)
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

func (b *Bot) GetMe() (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	apiResp, err := b.SendRequest(ctx, b.BaseUrl, RequestInfo{
		Method:      MethodGetMe,
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

	u := &types.User{}
	err = json.Unmarshal(apiResp.Result, u)
	return u, err
}

func (b *Bot) GetUpdates(params UpdateParams) ([]Update, error) {
	updates := []Update{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(params.Timeout))
	defer cancel()

	apiResp, err := b.SendRequest(ctx, b.BaseUrl, RequestInfo{
		Method:      MethodGetUpdates,
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

func (b *Bot) UploadFile(ctx context.Context, params UploadParams, file FileInfo) (*types.Message, error) {
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

	apiResp, err := b.SendRequest(ctx, b.BaseUrl, RequestInfo{
		Method:      "sendDocument",
		HasBody:     true,
		Body:        &WrappedReader{Reader: pipeReader},
		ContentType: multipartWriter.FormDataContentType(),
	})
	if err != nil {
		return nil, err
	}

	msg := &types.Message{}
	if err := json.NewDecoder(bytes.NewReader(apiResp.Result)).Decode(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (b *Bot) GetFile(fileId string) (*types.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	jbytes, err := json.Marshal(map[string]string{"file_id": fileId})
	if err != nil {
		return nil, err
	}

	apiResp, err := b.SendRequest(ctx, b.BaseUrl, RequestInfo{
		Method:      MethodGetFile,
		HasBody:     true,
		Body:        &WrappedReader{Reader: bytes.NewReader(jbytes)},
		ContentType: ContentTypeApplicationJson,
	})
	if err != nil {
		return nil, err
	}

	file := &types.File{}
	err = json.Unmarshal(apiResp.Result, file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (b *Bot) SendMessage(params TextMessageParams) (*types.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	apiResp, err := b.SendRequest(ctx, b.BaseUrl, RequestInfo{
		Method:      MethodSendMessage,
		HasBody:     true,
		Body:        &params,
		ContentType: params.ContentType(),
	})
	if err != nil {
		return nil, err
	}

	msg := &types.Message{}
	err = json.Unmarshal(apiResp.Result, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (b *Bot) EditMessageText(params EditMessageTextParams) (*types.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	apiResp, err := b.SendRequest(ctx, b.BaseUrl, RequestInfo{
		Method:      MethodEditMessageText,
		HasBody:     true,
		Body:        &params,
		ContentType: params.ContentType(),
	})
	if err != nil {
		return nil, err
	}

	msg := &types.Message{}
	err = json.Unmarshal(apiResp.Result, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (b *Bot) DeleteMessage(chatId, messageId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	jbytes, err := json.Marshal(map[string]int{"chat_id": chatId, "message_id": messageId})
	if err != nil {
		return err
	}

	_, err = b.SendRequest(ctx, b.BaseUrl, RequestInfo{
		Method:      MethodDeleteMessage,
		HasBody:     true,
		Body:        &WrappedReader{Reader: bytes.NewReader(jbytes)},
		ContentType: ContentTypeApplicationJson,
	})
	return err
}

func (b *Bot) StartPolling(params UpdateParams) (<-chan Update, error) {
	b.updatesChan = make(chan Update, params.Limit)

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
					update.bot = b
					b.updatesChan <- update
				}
			}
			time.Sleep(GetUpdatesSleepTime)
		}
	}()

	return b.updatesChan, nil
}
