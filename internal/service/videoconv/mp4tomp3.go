package videoconv

import (
	"bytes"
	"io"

	"bot/internal/core"
	"bot/lib/e"

	"github.com/kkdai/youtube/v2"
)

func DownloadSource(v *core.Video) error {
	err := downloadVideo(v)

	if err != nil {
		return e.Wrap("can't download video", err)
	}

	return nil
}

func ConvertToAudio() error {
	return nil
}

func downloadVideo(v *core.Video) error {
	videoStream, err := video(v.URL)
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

func video(url string) (io.ReadCloser, error) {
	client := youtube.Client{}

	videoInfo, err := client.GetVideo(url)

	if err != nil {
		return nil, e.Wrap("can't get video information", err)
	}

	stream, _, err := client.GetStream(videoInfo, &videoInfo.Formats[0])
	if err != nil {
		return nil, e.Wrap("can't get video stream", err)
	}

	return stream, nil
}
