package tgclient

import (
	"bot/internal/tech/coding"
	"bot/internal/tech/e"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (c *Client) ConfirmDeletionMessage(chatID, messageID int, title, username string) (err error) {
	defer func() { err = e.Wrap("can't confirm message deletion", err) }()

	hash := coding.EncodeUsernameAndTitle(username, title)

	replyMarkup, err := confirmDeletionMsgReplyMarkup(strconv.Itoa(messageID), hash)
	if err != nil {
		return err
	}

	url := c.baseURL(editMessageReplyMarkupMethod)

	jsonData, err := json.Marshal((map[string]interface{}{
		"chat_id":      chatID,
		"message_id":   messageID,
		"reply_markup": json.RawMessage(replyMarkup),
	}))
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

func (c *Client) RestoreDeletionMarkup(chatID, messageID int, hash string) (err error) {
	defer func() { err = e.Wrap("can't restore original keyboard", err) }()

	replyMarkup, err := deleteMsgReplyMarkup(hash)
	if err != nil {
		return err
	}

	url := c.baseURL(editMessageReplyMarkupMethod)

	jsonData, err := json.Marshal((map[string]interface{}{
		"chat_id":      chatID,
		"message_id":   messageID,
		"reply_markup": json.RawMessage(replyMarkup),
	}))
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

func confirmDeletionMsgReplyMarkup(messageID string, hash string) ([]byte, error) {
	confirmCallbackData := fmt.Sprintf("confirm_deletion:%s:%s", messageID, hash)
	refuseCallbackData := fmt.Sprintf("refuse_deletion:%s", hash)

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{
					"text":          "Да",
					"callback_data": confirmCallbackData,
				},
				{
					"text":          "Нет",
					"callback_data": refuseCallbackData,
				},
			},
		},
	}

	replyMarkupJSON, err := json.Marshal(replyMarkup)
	if err != nil {
		return nil, e.Wrap("can't create deletion message reply markup", err)
	}

	return replyMarkupJSON, nil
}
