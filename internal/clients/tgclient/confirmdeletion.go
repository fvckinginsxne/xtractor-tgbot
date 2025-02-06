package tgclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"bot/pkg/tech/coding"
	"bot/pkg/tech/e"
)

func (c *Client) ConfirmDeletionMessage(chatID, messageID int, title, username string) (err error) {
	defer func() { err = e.Wrap("can't confirm message deletion", err) }()

	hash := coding.EncodeUsernameAndTitle(username, title)

	replyMarkup, err := confirmDeletionMsgReplyMarkup(strconv.Itoa(messageID), hash)
	if err != nil {
		return err
	}

	url := c.baseURL(sendMessageMethod)

	confirmMsg := fmt.Sprintf("Удалить %s из плейлиста?", title)

	jsonData, err := json.Marshal((map[string]interface{}{
		"chat_id":      chatID,
		"text":         confirmMsg,
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

func confirmDeletionMsgReplyMarkup(messageID string, uuid string) ([]byte, error) {
	confirmCallbackData := fmt.Sprintf("confirm_deletion:%s:%s", messageID, uuid)

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{
					"text":          "Да",
					"callback_data": confirmCallbackData,
				},
				{
					"text":          "Нет",
					"callback_data": "refuse_deletion:",
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
