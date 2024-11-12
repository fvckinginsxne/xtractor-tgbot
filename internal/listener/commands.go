package listener

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"bot/internal/core"
	"bot/lib/e"
	"bot/pkg/storage"
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
		return l.savePage(text, chatID, username)
	}

	switch text {
	case RndCmd:
		return l.sendRandomPage(chatID, username)
	case HelpCmd:
		return l.sendHelp(chatID)
	case StartCmd:
		return l.sendGreeting(chatID)
	default:
		return l.tg.SendMessage(chatID, msgUnknownCmd)
	}
}

func (l *Listener) savePage(pageURL string, chatID int, username string) error {
	page := &core.Page{
		URL:      pageURL,
		Username: username,
	}

	isExists, err := l.storage.IsExists(page)
	if err != nil {
		return e.Wrap("can't save page", err)
	}

	if isExists {
		return l.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := l.storage.Save(page); err != nil {
		return e.Wrap("can't save page", err)
	}
	if err := l.tg.SendMessage(chatID, msgSaved); err != nil {
		return e.Wrap("can't save page", err)
	}

	return nil
}

func (l *Listener) sendRandomPage(chatID int, username string) error {
	page, err := l.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return e.Wrap("can't send random message", err)
	}

	if errors.Is(err, storage.ErrNoSavedPage) {
		return l.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := l.tg.SendMessage(chatID, page.URL); err != nil {
		return e.Wrap("can't send random page", err)
	}

	return l.storage.Remove(page)
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
