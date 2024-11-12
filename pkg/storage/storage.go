package storage

import (
	"errors"

	"bot/internal/core"
)

var ErrNoSavedPage = errors.New("no saved page")

type Storage interface {
	Save(p *core.Page) error
	PickRandom(username string) (*core.Page, error)
	Remove(p *core.Page) error
	IsExists(p *core.Page) (bool, error)
}
