package eventprocessor

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"bot/lib/e"
	"bot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *EventProcessor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from %s", text, username)

	if isAddCmd(text) {
		return p.savePage(text, chatID, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandomPage(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendGreeting(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCmd)
	}
}

func (p *EventProcessor) savePage(pageURL string, chatID int, username string) error {
	page := &storage.Page{
		URL:      pageURL,
		Username: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return e.Wrap("can't save page", err)
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return e.Wrap("can't save page", err)
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return e.Wrap("can't save page", err)
	}

	return nil
}

func (p *EventProcessor) sendRandomPage(chatID int, username string) error {
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return e.Wrap("can't send random message", err)
	}

	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return e.Wrap("can't send random page", err)
	}

	return p.storage.Remove(page)
}

func (p *EventProcessor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *EventProcessor) sendGreeting(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	url, err := url.Parse(text)

	return err == nil && url.Host != ""
}
