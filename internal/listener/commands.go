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
	HelpCmd  = "/help"
	StartCmd = "/start"
	ListCmd  = "/lst"
)

func (l *Listener) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from %s", text, username)

	if isAddCmd(text) {
		return l.saveAudio(text, chatID, username)
	}

	switch text {
	case HelpCmd:
		return l.sendHelp(chatID)
	case StartCmd:
		return l.sendGreeting(chatID)
	case ListCmd:
		return l.sendPlaylist(chatID, username)
	default:
		return l.tg.SendMessage(chatID, msgUnknownCmd)
	}
}

func (l *Listener) saveAudio(videoURL string, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't save video", err) }()

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

	if err := audio.DownloadData(); err != nil {
		return err
	}

	uuid := coding.EncodeUsernameAndTitle(username, audio.Title)

	if err := l.audioStorage.SaveAudio(audio, username, uuid); err != nil {
		return err
	}

	err = l.tg.SendAudio(chatID, audio.Data, audio.Title, username)
	if err != nil {
		return err
	}

	if err := l.tg.SendMessage(chatID, msgSaved); err != nil {
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

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	url, err := url.Parse(text)

	return err == nil && url.Host != ""
}
