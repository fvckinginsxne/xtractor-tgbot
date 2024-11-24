package core

import (
	"bytes"
	"io"

	"bot/lib/e"

	"github.com/kkdai/youtube/v2"
)

type Video struct {
	URL      string
	Source   []byte
	Username string
}

func (v *Video) DownloadSource() error {
	err := v.download()

	if err != nil {
		return e.Wrap("can't download video", err)
	}

	return nil
}

func ConvertToAudio() error {
	return nil
}

func (v *Video) download() error {
	videoStream, err := v.video()
	if err != nil {
		return e.Wrap("can't get video", err)
	}

	var buffer bytes.Buffer

	if _, err := io.Copy(&buffer, videoStream); err != nil {
		return e.Wrap("can't write data to buffer", err)
	}

	v.Source = buffer.Bytes()

	return nil
}

func (v *Video) video() (io.ReadCloser, error) {
	client := youtube.Client{}

	videoInfo, err := client.GetVideo(v.URL)
	if err != nil {
		return nil, e.Wrap("can't get video information", err)
	}

	stream, _, err := client.GetStream(videoInfo, &videoInfo.Formats[0])
	if err != nil {
		return nil, e.Wrap("can't get video stream", err)
	}

	return stream, nil
}
