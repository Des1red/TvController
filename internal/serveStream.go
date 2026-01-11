// internal/serverStream.go
package internal

import (
	"io"
	"net/http"
	"time"

	"tvctrl/logger"
)

func ServeStreamGo(
	cfg Config,
	stop <-chan struct{},
	streamPath string,
	container StreamContainer,
	source StreamSource,
) {
	cfg.ServerUp = true

	mux := http.NewServeMux()

	mux.HandleFunc(streamPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rc, err := source.Open()
		if err != nil {
			http.Error(w, "stream source unavailable", http.StatusServiceUnavailable)
			return
		}
		defer rc.Close()

		// w.Header().Set("Content-Type", container.ContentType())
		w.Header().Set("Content-Type", "video/mpeg")

		// For live-ish behavior, don't set Content-Length.
		// Many TVs accept chunked transfer; Go will do it automatically if no length is set.
		w.Header().Set("Accept-Ranges", "none")

		// HEAD: only headers
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Stream bytes
		// NOTE: This is phase-1 linear streaming. Range/seek comes later.
		_, _ = io.Copy(w, rc)
	})

	srv := &http.Server{
		Addr:              "0.0.0.0:" + cfg.ServePort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Success("Go HTTP stream server listening: %s%s (type=%s)", "http://"+cfg.LIP+":"+cfg.ServePort, streamPath, container.Key())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP stream server error: %v", err)
		}
	}()

	go func() {
		<-stop
		logger.Notify("Shutting down stream HTTP server")
		_ = srv.Close()
	}()
}
