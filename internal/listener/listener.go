package listener

import (
	"bot/internal/clients/tgclient"
	"bot/pkg/postgres/audiostorage"
	"bot/pkg/postgres/userstorage"
)

type Listener struct {
	tg           *tgclient.Client
	offset       int
	audioStorage *audiostorage.AudioStorage
	userStorage  *userstorage.UserStorage
}

func New(client *tgclient.Client, audioStorage *audiostorage.AudioStorage,
	userStorage *userstorage.UserStorage) *Listener {
	return &Listener{
		tg:           client,
		audioStorage: audioStorage,
		userStorage:  userStorage,
	}
}
