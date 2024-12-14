package tgclient

import "net/http"

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod         = "getUpdates"
	sendMessageMethod        = "sendMessage"
	sendAudioMethod          = "sendAudio"
	deleteMessageMethod      = "deleteMessage"
	sendCallbackAnswerMethod = "answerCallbackQuery"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}
