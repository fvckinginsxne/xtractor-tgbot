package storage

import (
	"errors"

	"bot/internal/core"
)

var ErrNoSavedPage = errors.New("no saved sources")

type Storage interface {
	Save(p *core.Audio) error
	Remove(p *core.Audio) error
	IsExists(p *core.Audio) (bool, error)
}
