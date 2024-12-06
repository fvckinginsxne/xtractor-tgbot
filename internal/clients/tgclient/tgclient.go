package tgclient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"bot/internal/core"
	"bot/lib/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
	sendAudioMethod   = "sendAudio"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

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

func (c *Client) SendAudio(chatID int, audio []byte, title string) (err error) {
	defer func() { err = e.Wrap("can't send audio", err) }()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("chat_id", strconv.Itoa(chatID)); err != nil {
		return err
	}

	if err := writer.WriteField("title", title); err != nil {
		return err
	}

	audioPart, err := writer.CreateFormFile("audio", "audio.mp3")
	if err != nil {
		return err
	}

	if _, err := audioPart.Write(audio); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	url := c.baseURL(sendAudioMethod)

	req, err := http.NewRequest(http.MethodPost, url.String(), &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Printf("Request URL: %s", req.URL.String())
	log.Println("Request Headers:")
	for key, values := range req.Header {
		for _, value := range values {
			log.Printf("%s: %s", key, value)
		}
	}
	log.Println("Buffer length:", buf.Len())

	resp, err := c.response(req)
	if err != nil {
		return err
	}

	if _, err := c.readResponse(resp); err != nil {
		log.Printf("Error reading response: %v", err)
		return err
	}

	log.Printf("Response Status: %s", resp.Status)

	return nil
}

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

func newBasePath(token string) string {
	return "bot" + token
}
