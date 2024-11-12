package listener

import (
	"errors"

	"bot/internal/clients/tgclient"
	"bot/internal/core"
	"bot/lib/e"
	"bot/pkg/storage"
)

type Listener struct {
	tg      *tgclient.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	chatID   int
	username string
}

var (
	ErrNoUpdates        = errors.New("no updates")
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *tgclient.Client, storage storage.Storage) *Listener {
	return &Listener{
		tg:      client,
		storage: storage,
	}
}

func (l *Listener) Fetch(limit int) ([]core.Event, error) {
	updates, err := l.tg.Updates(l.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get updates from telegram api", err)
	}

	if len(updates) == 0 {
		return nil, e.Wrap("no updates", ErrNoUpdates)
	}

	res := make([]core.Event, 0, len(updates))

	for _, update := range updates {
		res = append(res, event(update))
	}

	l.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (l *Listener) Process(event core.Event) error {
	switch event.Type {
	case core.Message:
		return l.processMessage(event)
	default:
		return e.Wrap("can't process unknown message", ErrUnknownEventType)
	}
}

func (l *Listener) processMessage(event core.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := l.doCmd(event.Text, meta.chatID, meta.username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event core.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(update core.Update) core.Event {
	updType := fetchType(update)

	res := core.Event{
		Type: updType,
		Text: fetchText(update),
	}

	if updType == core.Message {
		res.Meta = Meta{
			chatID:   update.Message.Chat.ID,
			username: update.Message.From.Username,
		}
	}

	return res
}

func fetchType(update core.Update) core.Type {
	if update.Message == nil {
		return core.Unknown
	}
	return core.Message
}

func fetchText(update core.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}
