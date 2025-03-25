package tgclient

import (
	"bot/internal/tech/coding"
	"bot/internal/tech/e"
	"bytes"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (c *Client) SendAudio(chatID int, audio []byte, title, username string) (err error) {
	defer func() { err = e.Wrap("can't send audio", err) }()

	hash := coding.EncodeUsernameAndTitle(username, title)

	deleteBtn, err := deleteMsgReplyMarkup(hash)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	err = addFields(writer, chatID, audio, title, deleteBtn)
	if err != nil {
		return err
	}

	url := c.baseURL(sendAudioMethod)

	req, err := http.NewRequest(http.MethodPost, url.String(), &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.response(req)
	if err != nil {
		return err
	}

	log.Printf("Response Status: %s", resp.Status)

	return nil
}

func deleteMsgReplyMarkup(hash string) (string, error) {
	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{
					"text":          "🗑",
					"callback_data": "delete_audio:" + hash,
				},
			},
		},
	}

	replyMarkupJSON, err := json.Marshal(replyMarkup)
	if err != nil {
		return "", e.Wrap("can't marshal replyMarkup to JSON", err)
	}

	log.Println("reply markup for message deletion: ", string(replyMarkupJSON))

	return string(replyMarkupJSON), nil
}

func addFields(writer *multipart.Writer, chatID int, audio []byte,
	title string, replyMarkup string) error {

	if err := writer.WriteField("chat_id", strconv.Itoa(chatID)); err != nil {
		return err
	}

	if err := writer.WriteField("title", title); err != nil {
		return err
	}

	if err := writer.WriteField("reply_markup", replyMarkup); err != nil {
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

	return nil
}
