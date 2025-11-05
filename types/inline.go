package types

type InlineQuery struct {
	Id       string
	From     User
	Query    string
	Offset   string
	ChatType string    `json:"chat_type,omitempty"`
	Location *Location `json:"location,omitempty"`
}

type ChosenInlineResult struct {
	ResultId        string `json:"result_id"`
	From            User
	Location        *Location `json:"location,omitempty"`
	InlineMessageId string    `json:"inline_message_id,omitempty"`
	Query           string
}

type CallbackQuery struct {
	Id              string
	From            User
	Message         *MaybeInaccessibleMessage `json:"message,omitempty"`
	InlineMessageId string                    `json:"inline_message_id,omitempty"`
	ChatInstance    string                    `json:"chat_instance"`
	Data            string
	GameShortName   string `json:"game_short_name,omitempty"`
}
