// internal/stream_play.go
package internal

import (
	"tvctrl/internal/avtransport"
)

func resolveAndPlayStream(cfg Config, streamPath string) {
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
