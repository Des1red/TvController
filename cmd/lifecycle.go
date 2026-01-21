package cmd

import (
	"os"
	"os/signal"
	"renderctl/internal/servers"
	"renderctl/internal/stream"
	"renderctl/internal/utils"
	"renderctl/logger"
	"time"
)

func preRun() (chan struct{}, bool) {
	stop := make(chan struct{})
	serverRunning := false

	// ---- PRE-RUN LOGIC ----
	mode := utils.NormalizeMode(cfg.Mode)
	if mode != "scan" && !cfg.ProbeOnly {
		inspectfile(mode)

		if mode != "stream" {
			servers.InitDefaultServer(cfg, stop)
		} else {
			stream.InitStreamServer(&cfg, stop)
		}

		time.Sleep(500 * time.Millisecond)
		serverRunning = true
	}

	return stop, serverRunning
}

func inspectfile(mode string) {
	if cfg.LFile == "" {
		logger.Error("Missing -Lf (media file)")
	}

	if mode != "stream" {
		if err := utils.ValidateFile(cfg.LFile); err != nil {
			logger.Error("Invalid file: %v", err)
		}
	}
}

func waitForShutdown(stop chan struct{}) {
	logger.Status("renderctl running — press Ctrl+C to exit")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	close(stop)
}
