package e

import (
	"errors"
	"fmt"
)

var (
	ErrNoUpdates          = errors.New("no updates")
	ErrUnknownEventType   = errors.New("unknown event type")
	ErrProcessTimedOut    = errors.New("yt-dlp process timed out")
	ErrLinkIsNotFromYT    = errors.New("link is not from youtube")
	ErrFileSizeIsTooLarge = errors.New("file size exceeds 50 MB")
)

func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
