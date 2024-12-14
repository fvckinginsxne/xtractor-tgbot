package tgclient

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"bot/internal/core"
	"bot/pkg/tech/e"
)

func (c *Client) Updates(offset, limit int) (updates []core.Update, err error) {
	defer func() { err = e.Wrap("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	url := c.baseURL(getUpdatesMethod)

	url.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.response(req)
	if err != nil {
		return nil, err
	}

	body, err := c.readResponse(resp)
	if err != nil {
		return nil, err
	}

	var res core.UpdatesResponse

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}
