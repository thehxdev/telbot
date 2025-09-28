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

	"github.com/thehxdev/telbot/types"
)

type Bot struct {
	token       string
	baseUrl     string
	baseFileUrl string
	self        *types.User
	client      http.Client
	updatesChan chan Update
}

type UpdateHandler func(update Update) error

type RequestInfo struct {
	Body        io.Reader
	Method      string
	ContentType string
}

type APIResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
}

func createMethodUrl(baseUrl string, method string) string {
	return fmt.Sprintf("%s/%s", baseUrl, method)
}

// Create a new Bot
func New(token string, host ...string) (*Bot, error) {
	h := "api.telegram.org"
	if len(host) > 0 {
		h = host[0]
	}

	b := &Bot{
		token:       token,
		baseUrl:     fmt.Sprintf("https://%s/bot%s", h, token),
		baseFileUrl: fmt.Sprintf("https://%s/file/bot%s", h, token),
		client:      http.Client{},
	}

	botUser, err := b.GetMe(context.Background())
	if err != nil {
		return nil, err
	}
	b.self = botUser

	return b, nil
}

func (b *Bot) SendRequest(ctx context.Context, baseUrl string, info RequestInfo) (resp APIResponse, err error) {
	var (
		reqBody  io.Reader
		httpResp *http.Response
	)

	if info.Body != nil {
		reqBody = info.Body
	}

	reqUrl := createMethodUrl(baseUrl, info.Method)
	req, err := http.NewRequestWithContext(ctx, "POST", reqUrl, reqBody)
	if err != nil {
		return
	}
	if info.ContentType != "" {
		req.Header.Set("Content-Type", info.ContentType)
	}

	httpResp, err = b.client.Do(req)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if err = json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return
	}

	if !resp.Ok {
		err = errors.New(resp.Description)
	}
	return
}

func (b *Bot) GetUpdates(ctx context.Context, params UpdateParams) ([]Update, error) {
	updates := []Update{}

	reqCtx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(params.Timeout))
	defer cancel()

	body, _ := ParamsToReader(params)
	apiResp, err := b.SendRequest(reqCtx, b.baseUrl, RequestInfo{
		Method:      MethodGetUpdates,
		Body:        body,
		ContentType: ContentTypeApplicationJson,
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

func (b *Bot) StartPolling(ctx context.Context, params UpdateParams) (<-chan Update, error) {
	b.updatesChan = make(chan Update, params.Limit)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			updates, err := b.GetUpdates(ctx, params)
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second * 5)
				continue
			}
			for _, update := range updates {
				if update.Id >= params.Offset {
					params.Offset = update.Id + 1
					update.Bot = b
					b.updatesChan <- update
				}
			}
			time.Sleep(getUpdatesSleepTime)
		}
	}()

	return b.updatesChan, nil
}

func (b *Bot) UploadFile(ctx context.Context, params UploadParams, files []IFileInfo) (*types.Message, error) {
	if len(files) == 0 {
		return nil, errors.New("no files provided to upload")
	}

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

		for _, file := range files {
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
				}
			}
		}
	}()

	apiResp, err := b.SendRequest(ctx, b.baseUrl, RequestInfo{
		Method:      "sendDocument",
		Body:        pipeReader,
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

func (b *Bot) GetMe(ctx context.Context) (*types.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultOperationTimeout)
	defer cancel()

	apiResp, err := b.SendRequest(ctx, b.baseUrl, RequestInfo{
		Method:      MethodGetMe,
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

func (b *Bot) LogOut(ctx context.Context) (bool, error) {
	reqCtx, cancel := context.WithTimeout(ctx, defaultOperationTimeout)
	defer cancel()
	apiResp, err := b.SendRequest(reqCtx, b.baseUrl, RequestInfo{
		Method: "logOut",
	})
	return apiResp.Ok, err
}

// This method is the implementation of the "close" method of telegram bot api
func (b *Bot) Close(ctx context.Context) (bool, error) {
	reqCtx, cancel := context.WithTimeout(ctx, defaultOperationTimeout)
	defer cancel()
	apiResp, err := b.SendRequest(reqCtx, b.baseUrl, RequestInfo{
		Method: "close",
	})
	return apiResp.Ok, err
}

func (b *Bot) GetFile(ctx context.Context, fileId string) (*types.File, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultOperationTimeout)
	defer cancel()

	body, _ := ParamsToReader(map[string]string{"file_id": fileId})
	apiResp, err := b.SendRequest(ctx, b.baseUrl, RequestInfo{
		Method:      MethodGetFile,
		Body:        body,
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

func (b *Bot) SendMessage(ctx context.Context, params TextMessageParams) (*types.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultOperationTimeout)
	defer cancel()

	body, _ := ParamsToReader(params)
	apiResp, err := b.SendRequest(ctx, b.baseUrl, RequestInfo{
		Method:      MethodSendMessage,
		Body:        body,
		ContentType: ContentTypeApplicationJson,
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

func (b *Bot) EditMessageText(ctx context.Context, params EditMessageTextParams) (*types.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultOperationTimeout)
	defer cancel()

	body, _ := ParamsToReader(params)
	apiResp, err := b.SendRequest(ctx, b.baseUrl, RequestInfo{
		Method:      MethodEditMessageText,
		Body:        body,
		ContentType: ContentTypeApplicationJson,
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

func (b *Bot) DeleteMessage(ctx context.Context, chatId, messageId int) error {
	ctx, cancel := context.WithTimeout(ctx, defaultOperationTimeout)
	defer cancel()

	body, _ := ParamsToReader(map[string]int{"chat_id": chatId, "message_id": messageId})
	_, err := b.SendRequest(ctx, b.baseUrl, RequestInfo{
		Method:      MethodDeleteMessage,
		Body:        body,
		ContentType: ContentTypeApplicationJson,
	})
	return err
}
