// internal/stream_play.go
package stream

import (
	"strings"
	"tvctrl/internal/avtransport"
	"tvctrl/internal/models"
	"tvctrl/internal/utils"
)

type StreamKind int

const (
	StreamFile StreamKind = iota
	StreamExternal
)

func ResolveStreamKind(cfg models.Config) StreamKind {
	if strings.HasPrefix(cfg.LFile, "http://") || strings.HasPrefix(cfg.LFile, "https://") {
		return StreamExternal
	}
	return StreamFile
}

func BuildStreamURL(cfg models.Config, streamPath string) string {
	p := strings.TrimPrefix(streamPath, "/")
	return "http://" + cfg.LIP + ":" + cfg.ServePort + "/" + p
}

func resolveAndPlayStream(cfg models.Config, streamPath string) {
	playStreamResolved(cfg, streamPath)
}

func playStreamResolved(cfg models.Config, streamPath string) {
	controlURL := utils.ControlURL(&cfg)
	if cfg.CachedControlURL != "" {
		controlURL = cfg.CachedControlURL
	}

	target := avtransport.Target{
		ControlURL: controlURL,
		MediaURL:   BuildStreamURL(cfg, streamPath),
	}

	meta := ""
	avtransport.Run(target, meta)
}
