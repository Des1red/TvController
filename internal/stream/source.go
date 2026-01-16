package stream

import (
	"errors"
	"tvctrl/internal/models"
	"tvctrl/internal/utils"
)

// This is generic in structure; later you replace with screen capture sources etc.
func BuildStreamSource(cfg models.Config) (StreamSource, error) {
	kind := ResolveStreamKind(cfg)

	switch kind {

	case StreamExternal:
		return urlSource{url: cfg.LFile}, nil

	case StreamFile:
		if err := utils.ValidateFile(cfg.LFile); err != nil {
			return nil, err
		}
		return fileSource{path: cfg.LFile}, nil
	}

	return nil, errors.New("unknown stream kind")
}
