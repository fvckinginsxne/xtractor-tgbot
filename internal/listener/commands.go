package listener

import (
	"log"
	"net/url"
	"strings"

	"bot/internal/core"
	"bot/lib/e"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (l *Listener) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from %s", text, username)

	if isAddCmd(text) {
		return l.saveVideo(text, chatID, username)
	}

	switch text {
	case HelpCmd:
		return l.sendHelp(chatID)
	case StartCmd:
		return l.sendGreeting(chatID)
	default:
		return l.tg.SendMessage(chatID, msgUnknownCmd)
	}
}

func (l *Listener) saveVideo(videoURL string, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't save video", err) }()

	audio := &core.Audio{
		URL:      videoURL,
		Username: username,
	}

	isExists, err := l.storage.IsExists(audio)
	if err != nil {
		return err
	}

	if isExists {
		return l.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := audio.DownloadSource(); err != nil {
		return err
	}

	if err := l.storage.Save(audio); err != nil {
		return err
	}

	if err := l.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (l *Listener) sendHelp(chatID int) error {
	return l.tg.SendMessage(chatID, msgHelp)
}

func (l *Listener) sendGreeting(chatID int) error {
	return l.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	url, err := url.Parse(text)

	return err == nil && url.Host != ""
}
