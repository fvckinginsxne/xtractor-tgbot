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
		return l.saveAudio(text, chatID, username)
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

func (l *Listener) saveAudio(videoURL string, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't save video", err) }()

	audio := &core.Audio{
		URL: videoURL,
	}

	isExists, err := l.audioStorage.IsExists(audio, username)
	if err != nil {
		return err
	}

	if isExists {
		return l.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := l.tg.SendMessage(chatID, msgProcessing); err != nil {
		return err
	}

	if err := audio.DownloadSource(); err != nil {
		return err
	}

	if err := l.audioStorage.Save(audio, username); err != nil {
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
