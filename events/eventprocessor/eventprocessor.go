package eventprocessor

import (
	"errors"

	"bot/clients/tgclient"
	"bot/events"
	"bot/lib/e"
	"bot/storage"
)

type EventProcessor struct {
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

func New(client *tgclient.Client, storage storage.Storage) *EventProcessor {
	return &EventProcessor{
		tg:      client,
		storage: storage,
	}
}

func (p *EventProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get updates from telegram api", err)
	}

	if len(updates) == 0 {
		return nil, e.Wrap("no updates", ErrNoUpdates)
	}

	res := make([]events.Event, 0, len(updates))

	for _, update := range updates {
		res = append(res, event(update))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *EventProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process unknown message", ErrUnknownEventType)
	}
}

func (p *EventProcessor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.chatID, meta.username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(update tgclient.Update) events.Event {
	updType := fetchType(update)

	res := events.Event{
		Type: updType,
		Text: fetchText(update),
	}

	if updType == events.Message {
		res.Meta = Meta{
			chatID:   update.Message.Chat.ID,
			username: update.Message.From.Username,
		}
	}

	return res
}

func fetchType(update tgclient.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(update tgclient.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}
