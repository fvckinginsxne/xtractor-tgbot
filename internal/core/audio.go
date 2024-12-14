package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"bot/pkg/tech/e"
)

type Audio struct {
	URL   string
	Data  []byte
	Title string
}

func (v *Audio) DownloadData() error {
	log.Printf("downloading...")

	err := v.download()
	if err != nil {
		return e.Wrap("can't download video", err)
	}

	log.Printf("audio is succesfully downloaded")

	return nil
}

func (v *Audio) download() error {
	outputFile := fmt.Sprintf("%s.mp3", "output")

	err := v.downloadVideoWithYTDLP(outputFile)
	if err != nil {
		return e.Wrap("failed to download video", err)
	}

	videoData, err := os.ReadFile(outputFile)
	if err != nil {
		return e.Wrap("failed to read downloaded video file", err)
	}

	v.Data = videoData

	defer os.Remove(outputFile)

	return nil
}

func (v *Audio) downloadVideoWithYTDLP(outputFile string) error {
	proxyURL := os.Getenv("PROXY_URL")

	log.Println(proxyURL)

	_, err := exec.LookPath("yt-dlp")
	if err != nil {
		return e.Wrap("yt-dlp not found", err)
	}

	titleCmd := exec.Command(
		"yt-dlp",
		"--print", "title",
		"--proxy", proxyURL,
		"--cookies-from-browser", "chrome",
		v.URL,
	)

	titleOutput, err := titleCmd.Output()
	if err != nil {
		return e.Wrap("failed to get video title", err)
	}

	v.Title = strings.TrimRight(string(titleOutput), "\n")

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
		v.URL,
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
			fmt.Printf("yt-dlp command failed: %v\n", err)
		}
	case <-time.After(30 * time.Second):
		fmt.Println("yt-dlp process timed out")
		downloadCmd.Process.Kill()
	}

	return nil
}
