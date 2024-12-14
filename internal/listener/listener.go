package listener

import (
	"bot/internal/clients/tgclient"
	"bot/pkg/sqlite/audiostorage"
	"bot/pkg/sqlite/userstorage"
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
