// internal/stream_mode.go
package internal

import (
	"errors"
	"strings"
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
	ServeStreamGo(cfg, stop, streamPath, container, src)

	// 5) Run discovery + AVTransport as usual, but pointing to /stream
	//    We do NOT use cfg.LFile/dir MediaURL here; we force stream URL.
	//cfg.StreamPath = streamPath                            // optional field if you add it; if not, ignore this line
	//cfg.OverrideMediaURL = BuildStreamURL(cfg, streamPath) // optional field if you add it; if not, use local var below

	// Reuse existing flow: runAuto/runManual ultimately call runWithConfig(cfg)
	// But cfg.MediaURL() currently returns your normal static URL.
	// Minimal, safe approach: bypass runWithConfig and construct target with stream URL directly:
	// controlURL := cfg.ControlURL()
	// if controlURL == "" && cfg.Mode == "stream" {
	// 	// Let existing auto/manual logic resolve it
	// }

	// We want your existing resolve logic. So:
	// - For stream, we still run auto/manual logic, but we need MediaURL() to be stream URL.
	// If you don't want to add fields to Config, do direct playback after resolve:

	// Resolve AVTransport endpoint using the same strategy as auto mode does
	// (we call runAuto but it uses cfg.MediaURL(). So we do the direct resolve here.)
	// Instead: copy the resolve logic minimal: SSDP -> cache -> probe, then play.

	// Minimal approach: runAuto-like resolution, then runWithTarget with stream URL:
	resolveAndPlayStream(cfg, container, streamPath)
}

// ---- helpers ----

func BuildStreamURL(cfg Config, streamPath string) string {
	p := strings.TrimPrefix(streamPath, "/")
	return "http://" + cfg.LIP + ":" + cfg.ServePort + "/" + p
}

// Container interface: designed to extend later (mp4, mkv, etc)
type StreamContainer interface {
	Key() string
	ContentType() string
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
func (tsContainer) ContentType() string {
	return "video/mpeg"
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
