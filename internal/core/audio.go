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
	URL      string
	Source   []byte
	Username string
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
	outputFile := fmt.Sprintf("%s.mp3", v.URL)

	err := v.downloadVideoWithYTDLP(outputFile)
	if err != nil {
		return e.Wrap("failed to download video", err)
	}

	videoData, err := os.ReadFile(outputFile)
	if err != nil {
		return e.Wrap("failed to read downloaded video file", err)
	}

	v.Source = videoData

	//defer os.Remove(outputFile)

	return nil
}

func (v *Audio) downloadVideoWithYTDLP(outputFile string) error {
	proxyURL := "http://DsGBKX:vyJezg@45.91.209.152:13162"

	_, err := exec.LookPath("yt-dlp")
	if err != nil {
		return e.Wrap("yt-dlp not found: make sure it is installed and in your PATH", err)
	}

	cmdArgs := []string{
		"-x",
		"--audio-format", "wav",
		"--proxy", proxyURL,
		"--cookies-from-browser", "chrome",
		"--no-post-overwrites",
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
	case <-time.After(30 * time.Second): // Устанавливаем таймаут 30 секунд
		fmt.Println("yt-dlp process timed out")
		cmd.Process.Kill() // Убиваем процесс, если он не завершился
	}

	// if err := cmd.Run(); err != nil {
	// 	return e.Wrap("yt-dlp command failed", err)
	// }

	return nil
}
