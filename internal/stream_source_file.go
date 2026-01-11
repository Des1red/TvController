package internal

import "os"

func (f fileSource) Open() (StreamReadCloser, error) {
	return os.Open(f.path)
}
