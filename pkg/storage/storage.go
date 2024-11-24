package storage

import (
	"errors"

	"bot/internal/core"
)

var ErrNoSavedPage = errors.New("no saved sources")

type Storage interface {
	Save(p *core.Video) error
	Remove(p *core.Video) error
	IsExists(p *core.Video) (bool, error)
}
