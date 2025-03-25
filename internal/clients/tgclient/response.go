package tgclient

import (
	"bot/internal/tech/e"
	"io"
	"net/http"
	"net/url"
	"path"
)

func (c *Client) response(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap("can't get response", err)
	}

	return resp, nil
}

func (c *Client) readResponse(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't read response", err)
	}

	_ = resp.Body.Close()

	return body, nil
}

func (c *Client) baseURL(method string) url.URL {
	return url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
}
