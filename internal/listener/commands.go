package listener

import (
	"log"
	"net/url"
	"strings"

	"bot/internal/core"
	"bot/pkg/tech/coding"
	"bot/pkg/tech/e"
)

const (
	HelpCmd     = "/help"
	StartCmd    = "/start"
	PlaylistCmd = "/lst"
)

func (l *Listener) doCmd(chatID int, text, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from %s", text, username)

	sentLink, err := isAddCmd(text)
	if err != nil {
		return l.tg.SendMessage(chatID, msgLinkIsNotFromYT)
	}

	if sentLink {
		return l.processVideoURL(chatID, text, username)
	}

	switch text {
	case HelpCmd:
		return l.sendHelp(chatID)
	case StartCmd:
		return l.sendGreeting(chatID)
	case PlaylistCmd:
		return l.sendPlaylist(chatID, username)
	default:
		return l.tg.SendMessage(chatID, msgUnknownCmd)
	}
}

func (l *Listener) processVideoURL(chatID int, videoURL, username string) (err error) {
	defer func() { err = e.Wrap("can't process video url", err) }()

	audio := &core.Audio{
		URL: videoURL,
	}

	isExists, err := l.audioStorage.IsExists(audio, username)
	if err != nil {
		return err
	}

	if isExists {
		return l.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := l.tg.SendMessage(chatID, msgProcessing); err != nil {
		return err
	}

	if err := l.saveAudio(audio, chatID, username); err != nil {
		return err
	}

	return nil
}

func (l *Listener) saveAudio(audio *core.Audio, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't save audio", err) }()

	if err := audio.ExtractAudio(); err != nil {
		return l.tg.SendMessage(chatID, msgErrorSavingAudio)
	}

	uuid := coding.EncodeUsernameAndTitle(username, audio.Title)

	if err := l.audioStorage.SaveAudio(audio, username, uuid); err != nil {
		return err
	}

	if err := l.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	err = l.tg.SendAudio(chatID, audio.Data, audio.Title, username)
	if err != nil {
		return err
	}

	return nil
}

func (l *Listener) sendPlaylist(chatID int, username string) error {
	audios, err := l.audioStorage.Playlist(username)
	if err != nil {
		return e.Wrap("can't get playlist", err)
	}

	if len(audios) == 0 {
		return l.tg.SendMessage(chatID, msgEmptyPlaylist)
	}

	for _, audio := range audios {
		log.Printf("Sending audio: Title=%s, DataSize=%d bytes", audio.Title, len(audio.Data))
		err := l.tg.SendAudio(chatID, audio.Data, audio.Title, username)
		if err != nil {
			return err
		}
		log.Printf("Successfully sent audio: %s", audio.Title)
	}

	return nil
}

func (l *Listener) sendHelp(chatID int) error {
	return l.tg.SendMessage(chatID, msgHelp)
}

func (l *Listener) sendGreeting(chatID int) error {
	return l.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) (bool, error) {
	return isURL(text)
}

func isURL(text string) (bool, error) {
	url, err := url.Parse(text)

	if err == nil && url.Host != "" {
		if isYTLink(text) {
			return true, nil
		} else {
			return false, e.ErrLinkIsNotFromYT
		}
	}

	return false, nil
}

func isYTLink(text string) bool {
	if strings.HasPrefix(text, "https://www.youtube.com/") {
		url, err := url.Parse(text)
		if err != nil {
			return false
		}

		q := url.Query()
		if videoID := q.Get("v"); videoID != "" {
			return true
		}
	}

	return false
}
