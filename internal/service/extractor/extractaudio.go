package extractor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"bot/internal/tech/e"
)

type Audio struct {
	AudioFile []byte
	Title     string
}

func ExtractAudio(videoURL string) (*Audio, error) {
	log.Printf("downloading audio...")

	audio := Audio{}

	proxyURL := os.Getenv("PROXY_URL")

	fmt.Println(proxyURL)

	title, err := videoTitle(proxyURL, videoURL)
	if err != nil {
		return nil, e.Wrap("can't download audio with yt-dlp", err)
	}

	audio.Title = title

	outputFile := fmt.Sprintf("%s.mp3", title)

	err = downloadAudioWithYTDLP(outputFile, videoURL, proxyURL)
	if err != nil {
		return nil, e.Wrap("failed to download audio", err)
	}

	ouputData, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, e.Wrap("failed to read downloaded audio file", err)
	}
	defer os.Remove(outputFile)

	if err := checkFileSize(outputFile); err != nil {
		return nil, e.Wrap("file size is too large", err)
	}

	audio.AudioFile = ouputData

	log.Printf("audio was succesfully downloaded")

	return &audio, nil
}

func downloadAudioWithYTDLP(outputFile, videoURL, proxyURL string) error {
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return e.Wrap("yt-dlp not found", err)
	}

	if err := updateYTDLP(); err != nil {
		return e.Wrap("can't download audio with yt-dlp", err)
	}

	if err := downloadAudioToFile(outputFile, videoURL, proxyURL); err != nil {
		return e.Wrap("can't download audio with yt-dlp", err)
	}

	return nil
}

func videoTitle(proxyURL, videoURL string) (string, error) {
	titleCmd := exec.Command(
		"yt-dlp",
		"--print", "title",
		"--proxy", proxyURL,
		"--cookies-from-browser", "chrome",
		videoURL,
	)

	titleOutput, err := titleCmd.Output()
	if err != nil {
		return "", e.Wrap("failed to get video title", err)
	}

	return strings.TrimRight(string(titleOutput), "\n"), nil
}

func downloadAudioToFile(outputFile, videoURL, proxyURL string) error {
	cmdArgs := []string{
		"-x",
		"--audio-format", "mp3",
		"--proxy", proxyURL,
		"--cookies-from-browser", "chrome",
		"--no-post-overwrites",
		"--retries", "10",
		"--fragment-retries", "10",
		"--socket-timeout", "30",
		"-o", outputFile,
		videoURL,
	}

	downloadCmd := exec.Command("yt-dlp", cmdArgs...)

	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr

	done := make(chan error)
	go func() {
		done <- downloadCmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return e.Wrap("yt-dlp command failed", err)
		}
	case <-time.After(60 * time.Second):
		downloadCmd.Process.Kill()
		return e.Wrap("yt-dlp process timed out", e.ErrProcessTimedOut)
	}

	return nil
}

func checkFileSize(outputFile string) error {
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		return e.Wrap("can't get file info", err)
	}

	const maxSize = 50 * 1024 * 1024
	if fileInfo.Size() > maxSize {
		return e.ErrFileSizeIsTooLarge
	}

	return nil
}

func updateYTDLP() error {
	updateCmd := exec.Command("yt-dlp", "-U")

	if err := updateCmd.Run(); err != nil {
		return e.Wrap("failed to update yt-dlp", err)
	}

	return nil
}
