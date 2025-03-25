package tgclient

import (
	"bot/internal/tech/e"
	"bytes"
	"encoding/json"
	"net/http"
)

func (c *Client) DeleteMessage(chatID int, messageID int) (err error) {
	defer func() { err = e.Wrap("can't delete message", err) }()

	url := c.baseURL(deleteMessageMethod)

	values := map[string]int{
		"chat_id":    chatID,
		"message_id": messageID,
	}

	jsonData, err := json.Marshal(values)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest(http.MethodPost, url.String(), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if _, err := c.response(req); err != nil {
		return err
	}

	return nil
}
