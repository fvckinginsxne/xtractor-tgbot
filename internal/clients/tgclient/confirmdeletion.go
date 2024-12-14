package tgclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"bot/pkg/tech/coding"
	"bot/pkg/tech/e"
	"bot/pkg/tech/marshaling"
)

func (c *Client) ConfirmDeletionMessage(chatID, messageID int, title, username string) (err error) {
	defer func() { err = e.Wrap("can't confirm message deletion", err) }()

	uuid := coding.EncodeUsernameAndTitle(username, title)

	replyMarkup, err := confirmDeletionMsgReplyMarkup(strconv.Itoa(messageID), uuid)
	if err != nil {
		return err
	}

	url := c.baseURL(sendMessageMethod)

	confirmMsg := fmt.Sprintf("Удалить %s из плейлиста?", title)

	reqBody, err := marshaling.DataToJSON((map[string]interface{}{
		"chat_id":      chatID,
		"text":         confirmMsg,
		"reply_markup": json.RawMessage(replyMarkup),
	}))
	if err != nil {
		return err
	}

	resp, err := http.Post(url.String(), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

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

	replyMarkupJSON, err := marshaling.DataToJSON(replyMarkup)
	if err != nil {
		return nil, e.Wrap("can't create deleton message reply markup", err)
	}

	log.Println("reply markup for confirm deletion: ", string(replyMarkupJSON))

	return replyMarkupJSON, nil
}
