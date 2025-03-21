package telbot

const (
	ChatTypePrivate    string = "private"
	ChatTypeGroup      string = "group"
	ChatTypeSuperGroup string = "supergroup"
	ChatTypeChannel    string = "channel"
)

const (
	ContentTypeFormUrlEncoded    = "application/x-www-form-urlencoded"
	ContentTypeMultipartFormData = "multipart/form-data"
	ContentTypeApplicationJson   = "application/json"
)

const (
	MethodGetMe           = "getMe"
	MethodGetUpdates      = "getUpdates"
	MethodSendMessage     = "sendMessage"
	MethodGetFile         = "getFile"
	MethodEditMessageText = "editMessageText"
	MethodDeleteMessage   = "deleteMessage"
)
