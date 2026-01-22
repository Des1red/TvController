package cmd

import (
	"os"

	"renderctl/internal"
	"renderctl/internal/cache"
	"renderctl/internal/models"
	"renderctl/internal/ui"
	"renderctl/internal/utils"
	"renderctl/logger"
	"renderctl/requirements"
)

var cfg = models.DefaultConfig
var noCache bool

func Execute() {
	parseFlags()
	handleInstaller()
	handleFlagsAndLogging()
	handleInteraction()

	stop, serverRunning := preRun()
	internal.RunScript(&cfg)

	if serverRunning {
		waitForShutdown(stop)
	}
	logger.CreateReport()
}

func handleInstaller() {
	// ---- INSTALLER (early exit) ----
	if requirements.Install {
		if err := requirements.RunInstaller(); err != nil {
			logger.Error("%v", err)
		}
		os.Exit(0)
	}
}

func handleFlagsAndLogging() {
	if bad, msg := badFlagUse(); bad {
		logger.Error(msg)
	}
	// TUI mode
	if cfg.Interactive {
		ui.Run(&cfg)
	}
	// Set verbose + filename(if any)
	logger.SetVerbose(cfg.Verbose)

	if cfg.ReportFileName != "" {
		cfg.ReportFile = true
		logger.SetFileName(cfg.ReportFileName)
	}
	// FLAG INVERSION
	cfg.UseCache = !noCache

	// Cache commands exit early
	if cache.HandleCacheCommands(cfg) {
		os.Exit(0)
	}

	if cfg.SelectCache >= 0 {
		cache.LoadCachedTV(&cfg)
	}
}

func handleInteraction() {
	if (cfg.Mode == "scan" && cfg.Discover) || (cfg.Mode != "scan" && !cfg.ProbeOnly) {
		cfg.LIP = utils.LocalIP(cfg.LIP)
	}
}
