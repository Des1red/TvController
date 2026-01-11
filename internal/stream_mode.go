// internal/stream_mode.go
package internal

import (
	"errors"
	"strings"
	"tvctrl/internal/avtransport"
	"tvctrl/logger"
)

func StartStreamMode(cfg Config, stop <-chan struct{}) {
	// 1) Resolve IP used for MediaURL
	cfg.LIP = LocalIP(cfg.LIP)

	// 2) Decide container (default TS for phase-1)
	containerKey := "ts"
	container, err := GetContainer(containerKey)
	if err != nil {
		logger.Fatal("Stream container error: %v", err)
		return
	}

	// 3) Decide source (phase-1: from -Lf if provided; later: screen / remote / etc)
	src, err := BuildStreamSource(cfg)
	if err != nil {
		logger.Fatal("Stream source error: %v", err)
		return
	}

	// 4) Start stream HTTP server (/stream)
	streamPath := "/stream"
	// Resolve AVTransport ONCE
	if cfg.UseCache {
		if tryCache(&cfg) {
			goto resolved
		}
	}

	if !tryProbe(&cfg) {
		logger.Fatal("Unable to resolve AVTransport endpoint")
	}

resolved:

	controlURL := cfg.ControlURL()
	if controlURL == "" {
		logger.Fatal("ControlURL not resolved")
	}

	// Now fetch protocol info from the SAME ControlURL
	media, err := avtransport.FetchMediaProtocols(controlURL)
	if err != nil {
		logger.Notify("ProtocolInfo fetch failed, using fallback MIME")
	}

	// choose MIME
	mime := selectMime(container, media)
	logger.Notify("Using stream MIME: %s", mime)
	ServeStreamGo(cfg, stop, streamPath, mime, src)
	resolveAndPlayStream(cfg, streamPath)
}

// ---- helpers ----

func BuildStreamURL(cfg Config, streamPath string) string {
	p := strings.TrimPrefix(streamPath, "/")
	return "http://" + cfg.LIP + ":" + cfg.ServePort + "/" + p
}

// Container interface: designed to extend later (mp4, mkv, etc)
type StreamContainer interface {
	Key() string
	MimeCandidates() []string
}

// Source interface: later you can implement screen capture, remote proxy, etc
type StreamSource interface {
	Open() (StreamReadCloser, error)
}

type StreamReadCloser interface {
	Read(p []byte) (int, error)
	Close() error
}

// Registry for containers
var containerRegistry = map[string]StreamContainer{
	"ts": tsContainer{},
}

func GetContainer(key string) (StreamContainer, error) {
	c, ok := containerRegistry[strings.ToLower(strings.TrimSpace(key))]
	if !ok {
		return nil, errors.New("unknown container: " + key)
	}
	return c, nil
}

type tsContainer struct{}

func (tsContainer) Key() string { return "ts" }

func (tsContainer) MimeCandidates() []string {
	return []string{
		"video/mpeg",               // most DLNA TVs accept this
		"application/octet-stream", // very permissive fallback
		"video/mp2t",               // least compatible
	}
}

// Phase-1 source builder: if -Lf is provided, stream that file’s bytes.
// This is generic in structure; later you replace with screen capture sources etc.
func BuildStreamSource(cfg Config) (StreamSource, error) {
	if strings.TrimSpace(cfg.LFile) == "" {
		return nil, errors.New("stream mode currently requires -Lf as a byte source (later: screen/camera/remote sources will be added)")
	}
	if err := ValidateFile(cfg.LFile); err != nil {
		return nil, err
	}
	return fileSource{path: cfg.LFile}, nil
}

type fileSource struct{ path string }

// implemented in internal/stream_source_file.go

func selectMime(
	container StreamContainer,
	supported map[string][]string,
) string {

	// If TV returned nothing → fallback
	if len(supported) == 0 {
		return "video/mpeg"
	}

	for _, cand := range container.MimeCandidates() {
		if _, ok := supported[cand]; ok {
			return cand
		}
	}

	// Nothing matched → safe fallback
	return "video/mpeg"
}
