package listener

import (
	"log"
	"strconv"
	"strings"

	"bot/internal/core"
	"bot/pkg/tech/e"
)

const deleteAudioPrefix = "delete_audio:"
const confirmDeletionAudioPrefix = "confirm_deletion:"
const refuseDeletionAudioPrefix = "refuse_deletion:"

func (l *Listener) Process(event core.Event) error {
	switch event.Type {
	case core.Message:
		return l.processMessage(event)
	case core.Data:
		return l.processCallback(event)
	default:
		return e.Wrap("can't process unknown message", e.ErrUnknownEventType)
	}
}

func (l *Listener) processCallback(event core.Event) (err error) {
	defer func() { err = e.Wrap("can't process callback request", err) }()

	chatID := event.ChatID
	messageID := event.MessageID

	data := event.CallbackQuery.Data

	if strings.HasPrefix(data, deleteAudioPrefix) {
		return l.processDeleteMsgCallback(event, data, chatID, messageID)
	} else if strings.HasPrefix(data, confirmDeletionAudioPrefix) {
		return l.processConfirmDeletionCallback(event, data, chatID, messageID)
	} else if strings.HasPrefix(data, refuseDeletionAudioPrefix) {
		return l.processRefuseDeletionCallback(event, chatID, messageID)
	}

	return nil
}

func (l *Listener) processDeleteMsgCallback(event core.Event, data string,
	chatID, messageID int) (err error) {
	defer func() { err = e.Wrap("can't process message deletion", err) }()

	title, username, err := l.parseData(data, deleteAudioPrefix)
	if err != nil {
		return err
	}

	log.Printf("Title=%s, username=%s", title, username)

	if err := l.tg.ConfirmDeletionMessage(chatID, messageID, title, username); err != nil {
		return err
	}

	return l.tg.SendCallback(event.CallbackQuery.ID)
}

func (l *Listener) processConfirmDeletionCallback(event core.Event, data string,
	chatID, messageID int) (err error) {
	defer func() { err = e.Wrap("can't process deletion confirmation", err) }()

	log.Println("deleting audio...")

	title, _, err := l.parseData(data, confirmDeletionAudioPrefix)
	if err != nil {
		return err
	}

	if err := l.tg.DeleteMessage(chatID, messageID); err != nil {
		return err
	}

	parsedMsgID, err := parseMsgID(data)
	if err != nil {
		return err
	}

	if err := l.tg.DeleteMessage(chatID, parsedMsgID); err != nil {
		return err
	}

	if err := l.audioStorage.Remove(title); err != nil {
		return err
	}

	log.Println("audio have successfully deleted")

	return l.tg.SendCallback(event.CallbackQuery.ID)
}

func (l *Listener) processRefuseDeletionCallback(event core.Event, chatID, messageID int) error {
	data := event.CallbackQuery.Data

	log.Print(data)

	hash := parseHash(data, refuseDeletionAudioPrefix)

	log.Print(hash)

	if err := l.tg.RestoreDeletionMarkup(chatID, messageID, hash); err != nil {
		return e.Wrap("can't process refuse deletion", err)
	}

	return l.tg.SendCallback(event.CallbackQuery.ID)
}

func (l *Listener) processMessage(event core.Event) error {
	if err := l.doCmd(event.ChatID, event.Text, event.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func (l *Listener) parseData(data, prefix string) (title, username string, err error) {
	hash := parseHash(data, prefix)

	title, err = l.audioStorage.Title(hash)
	if err != nil {
		return "", "", e.Wrap("can't parse data", err)
	}

	username, err = l.userStorage.Username(hash)
	if err != nil {
		return "", "", e.Wrap("can't parse data", err)
	}

	return title, username, nil
}

func parseMsgID(callbackData string) (int, error) {
	parts := strings.Split(callbackData, ":")

	msgID, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, e.Wrap("can't parse message id", err)
	}

	return msgID, nil
}

func parseHash(calbackData, prefix string) string {
	var parts []string
	if prefix == confirmDeletionAudioPrefix {
		parts = strings.Split(calbackData, ":")
		return parts[2]
	}

	hash := strings.TrimPrefix(calbackData, prefix)

	return hash
}
