// internal/stream_play.go
package internal

import (
	"tvctrl/internal/avtransport"
	"tvctrl/logger"
)

func resolveAndPlayStream(cfg Config, container StreamContainer, streamPath string) {
	// Resolve AVTransport endpoint the same way "auto" does
	// // 1) SSDP
	// if trySSDP(&cfg) {
	// 	playStreamResolved(cfg, streamPath)
	// 	return
	// }

	// 2) Cache
	if cfg.UseCache {
		if tryCache(&cfg) {
			playStreamResolved(cfg, streamPath)
			return
		}
	}

	// 3) Probe fallback
	ok := tryProbe(&cfg)
	if !ok {
		logger.Fatal("Unable to resolve AVTransport endpoint")
		return
	}

	if cfg.ProbeOnly {
		logger.Success("Probe completed (no playback).")
		return
	}

	playStreamResolved(cfg, streamPath)
}

func playStreamResolved(cfg Config, streamPath string) {
	controlURL := cfg.ControlURL()
	if cfg._CachedControlURL != "" {
		controlURL = cfg._CachedControlURL
	}

	target := avtransport.Target{
		ControlURL: controlURL,
		MediaURL:   BuildStreamURL(cfg, streamPath),
	}

	meta := ""
	avtransport.Run(target, meta)
}
