package tgclient

import (
	"bot/internal/tech/e"
	"bytes"
	"encoding/json"
	"net/http"
)

func (c *Client) SetCommandsList() (err error) {
	defer func() { err = e.Wrap("can't set commands list", err) }()

	commands := map[string]interface{}{
		"commands": []map[string]string{
			{
				"command":     "/start",
				"description": "начать работу с ботом",
			},
			{
				"command":     "/help",
				"description": "получить описание работы с ботом",
			},
			{
				"command":     "/lst",
				"description": "получить плейлист",
			},
		},
	}

	jsonData, err := json.Marshal(commands)
	if err != nil {
		return err
	}

	url := c.baseURL(setCommandsListMethod)

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
