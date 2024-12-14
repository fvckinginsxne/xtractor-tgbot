package core

type Type int

const (
	Unknown Type = iota
	Message
	Data
)

type Event struct {
	Type          Type
	Text          string
	MessageID     int
	ChatID        int
	Username      string
	CallbackQuery *CallbackQuery
}
