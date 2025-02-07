package tgclient

import (
	"bytes"
	"encoding/json"
	"net/http"

	"bot/pkg/tech/e"
)

func (c *Client) SendCallback(callbackID string) (err error) {
	defer func() { err = e.Wrap("can't send callback answer", err) }()

	url := c.baseURL(sendCallbackAnswerMethod)

	data := map[string]string{
		"callback_query_id": callbackID,
	}

	jsonData, err := json.Marshal(data)
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
