package tgclient

import (
	"net/http"
	"net/url"
	"strconv"

	"bot/pkg/tech/e"
)

func (c *Client) SendMessage(chatID int, text string) (err error) {
	defer func() { err = e.Wrap("can't send message", err) }()

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	url := c.baseURL(sendMessageMethod)

	url.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}

	if _, err := c.response(req); err != nil {
		return err
	}

	return nil
}
