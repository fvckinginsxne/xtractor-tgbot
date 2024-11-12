package core

import (
	"crypto/sha1"
	"fmt"
	"io"

	"bot/lib/e"
)

type Page struct {
	URL      string
	Username string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("cant't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.Username); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
