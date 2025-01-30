package listener

import (
	"bot/internal/core"
	"bot/pkg/tech/e"
)

func (l *Listener) Fetch(limit int) ([]core.Event, error) {
	updates, err := l.tg.Updates(l.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get updates from telegram api", err)
	}

	if len(updates) == 0 {
		return nil, e.Wrap("no updates", e.ErrNoUpdates)
	}

	res := make([]core.Event, 0, len(updates))

	for _, update := range updates {
		res = append(res, event(update))
	}

	l.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func event(update core.Update) core.Event {
	updType := fetchType(update)

	res := core.Event{
		Type:          updType,
		Text:          fetchText(update),
		CallbackQuery: update.CallbackQuery,
	}

	if update.CallbackQuery != nil {
		res.ChatID = update.CallbackQuery.Message.Chat.ID
		res.MessageID = update.CallbackQuery.Message.MessageID
	}

	if update.Message != nil {
		res.MessageID = update.Message.MessageID
		res.ChatID = update.Message.Chat.ID
		res.Username = update.Message.From.Username
	}

	return res
}

func fetchType(update core.Update) core.Type {
	if update.CallbackQuery != nil {
		return core.Data
	}

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
