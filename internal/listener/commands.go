package listener

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"bot/internal/service/extractor"
	"bot/internal/tech/coding"
	"bot/internal/tech/e"
)

const (
	HelpCmd     = "/help"
	StartCmd    = "/start"
	PlaylistCmd = "/lst"
)

func (l *Listener) doCmd(chatID int, text, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from %s", text, username)

	sentLink, err := isURL(text)
	if err == e.ErrLinkIsNotFromYT {
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

	isExists, err := l.audioStorage.IsExists(videoURL)
	if err != nil {
		return err
	}

	if isExists {
		return l.tg.SendMessage(chatID, msgAudioAlreadyExists)
	}

	if err := l.tg.SendMessage(chatID, msgProcessing); err != nil {
		return err
	}

	audio, err := extractor.ExtractAudio(videoURL)
	if errors.Is(err, e.ErrFileSizeIsTooLarge) {
		return l.tg.SendMessage(chatID, msgFileIsToLarge)
	} else if err != nil {
		return l.tg.SendMessage(chatID, msgErrorSavingAudio)
	}

	if err := l.saveAudio(audio, chatID, videoURL, username); err != nil {
		return err
	}

	return nil
}

func (l *Listener) saveAudio(audio *extractor.Audio, chatID int,
	videoURL, username string) (err error) {
	defer func() { err = e.Wrap("can't save audio", err) }()

	hash := coding.EncodeUsernameAndTitle(username, audio.Title)

	userID, err := l.userStorage.Save(username)
	if err != nil {
		return err
	}

	urlID, err := l.urlstorage.Save(videoURL)
	if err != nil {
		return err
	}

	if err := l.audioStorage.Save(audio, hash, userID, urlID); err != nil {
		return err
	}

	if err := l.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	err = l.tg.SendAudio(chatID, audio.AudioFile, audio.Title, username)
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
		log.Printf("Sending audio: Title=%s", audio.Title)
		err := l.tg.SendAudio(chatID, audio.AudioFile, audio.Title, username)
		if err != nil {
			return e.Wrap("can't send playlist", err)
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

func isURL(text string) (bool, error) {
	url, _ := url.Parse(text)

	if url.Host != "" {
		if isYT(url) {
			return true, nil
		} else {
			return false, e.ErrLinkIsNotFromYT
		}
	}

	return false, nil
}

func isYT(url *url.URL) bool {
	if url.Host == "www.youtube.com" || url.Host == "youtube.com" {
		return url.Path == "/watch" && url.Query().Get("v") != ""
	}

	if url.Host == "youtu.be" {
		return url.Path != "" && len(strings.Trim(url.Path, "/")) > 0
	}

	return false
}
