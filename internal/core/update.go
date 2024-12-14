package core

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message"`
	CallbackQuery *CallbackQuery   `json:"callback_query,omitempty"`
}

type IncomingMessage struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
}

type CallbackQuery struct {
	ID      string           `json:"id"`
	Data    string           `json:"data,omitempty"`
	Message *IncomingMessage `json:"message,omitempty"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
