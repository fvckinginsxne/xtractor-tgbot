package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"bot/lib/e"
)

type Audio struct {
	URL   string
	Data  []byte
	Title string
}

func (v *Audio) DownloadSource() error {
	log.Printf("downloading...")

	err := v.download()

	if err != nil {
		return e.Wrap("can't download video", err)
	}

	log.Printf("audio is succesfully downloaded")

	return nil
}

func (v *Audio) download() error {
	outputFile := fmt.Sprintf("%s.wav", "output")

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

	v.Title = string(titleOutput)

	cmdArgs := []string{
		"-x",
		"--audio-format", "wav",
		"--proxy", proxyURL,
		"--cookies-from-browser", "chrome",
		"--no-post-overwrites",
		"--retries", "10", // Повторить попытку до 10 раз
		"--fragment-retries", "10", // Повторять попытку при сбое фрагмента
		"--socket-timeout", "30",
		"-o", outputFile,
		v.URL,
	}

	cmd := exec.Command("yt-dlp", cmdArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	done := make(chan error)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("yt-dlp command failed: %v\n", err)
		}
	case <-time.After(30 * time.Second):
		fmt.Println("yt-dlp process timed out")
		cmd.Process.Kill()
	}

	return nil
}
