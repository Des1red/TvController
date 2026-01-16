package stream

import (
	"errors"
	"strings"
)

type tsContainer struct{}

func (tsContainer) Key() string { return "ts" }

func (tsContainer) MimeCandidates() []string {
	return []string{
		"video/mpeg",               // most DLNA TVs accept this
		"application/octet-stream", // very permissive fallback
		"video/mp2t",               // least compatible
	}
}

type passthroughContainer struct{}

func (p passthroughContainer) Key() string { return "passthrough" }

func (p passthroughContainer) MimeCandidates() []string {
	return []string{
		"video/mp4",
		"video/mpeg",
		"application/octet-stream",
	}
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
	"ts":          tsContainer{},
	"mp4":         mp4Container{},
	"passthrough": passthroughContainer{},
}

type mp4Container struct{}

func (mp4Container) Key() string { return "mp4" }

func (mp4Container) MimeCandidates() []string {
	return []string{
		"video/mp4",
		"application/octet-stream",
	}
}

func GetContainer(key string) (StreamContainer, error) {
	c, ok := containerRegistry[strings.ToLower(strings.TrimSpace(key))]
	if !ok {
		return nil, errors.New("unknown container: " + key)
	}
	return c, nil
}
