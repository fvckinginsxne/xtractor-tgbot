package listener

import (
	"bot/internal/clients/tgclient"
	"bot/internal/postgres/audiostorage"
	"bot/internal/postgres/urlstorage"
	"bot/internal/postgres/userstorage"
)

type Listener struct {
	tg           *tgclient.Client
	offset       int
	audioStorage *audiostorage.AudioStorage
	userStorage  *userstorage.UserStorage
	urlstorage   *urlstorage.UrlStorage
}

func New(client *tgclient.Client, audioStorage *audiostorage.AudioStorage,
	userStorage *userstorage.UserStorage, urlStorage *urlstorage.UrlStorage) *Listener {
	return &Listener{
		tg:           client,
		audioStorage: audioStorage,
		userStorage:  userStorage,
		urlstorage:   urlStorage,
	}
}
