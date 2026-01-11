package cmd

import (
	"os"
	"os/signal"
	"time"

	"tvctrl/internal"
	"tvctrl/logger"
)

var cfg = internal.DefaultConfig
var noCache bool

func Execute() {
	parseFlags()

	// FLAG INVERSION
	cfg.UseCache = !noCache

	// Cache commands exit early
	if internal.HandleCacheCommands(cfg) {
		os.Exit(0)
	}

	if cfg.SelectCache >= 0 {
		internal.LoadCachedTV(&cfg)
	}

	stop := make(chan struct{})
	serverRunning := false

	// ---- PRE-RUN LOGIC ----
	mode := internal.NormalizeMode(cfg.Mode)
	if mode != "scan" && mode != "stream" && !cfg.ProbeOnly {
		if cfg.LFile == "" {
			logger.Fatal("Missing -Lf (media file)")
			os.Exit(1)
		}

		if err := internal.ValidateFile(cfg.LFile); err != nil {
			logger.Fatal("Invalid file: %v", err)
			os.Exit(1)
		}

		cfg.LIP = internal.LocalIP(cfg.LIP)
		internal.ServeDirGo(cfg, stop)
		time.Sleep(500 * time.Millisecond)
		serverRunning = true
	}
	// Stream mode starts its own HTTP server (so keep the process alive)
	if mode == "stream" && !cfg.ProbeOnly {
		serverRunning = true
	}
	if mode == "stream" {
		logger.Notify("Waiting 10 seconds before starting playback...")
		time.Sleep(10 * time.Second)
	}

	internal.RunScript(cfg, stop)

	if !serverRunning {
		return
	}

	logger.Info("tvctrl running — press Ctrl+C to exit")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	close(stop)
}
