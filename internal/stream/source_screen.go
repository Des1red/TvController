package stream

import (
	"os"
	"os/exec"
	"tvctrl/logger"
)

func newRollingFileSource() *rollingFileSource {
	path := "/tmp/screen_fifo"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := exec.Command("mkfifo", path).Run(); err != nil {
			logger.Fatal("Failed to create FIFO: %v", err)
		}
	}

	cmd := exec.Command("ffmpeg",
		"-f", "x11grab",
		"-i", ":0.0",
		"-r", "25",
		"-f", "mp4",
		"-movflags", "frag_keyframe+empty_moov",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-tune", "zerolatency",
		"-profile:v", "baseline",
		"-level", "3.0",
		"-b:v", "1M",
		"-y",
		path,
	)

	return &rollingFileSource{
		path:    path,
		cmd:     cmd,
		started: false,
		opened:  false,
	}
}
