package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, wroteHeader: false}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.ResponseWriter.WriteHeader(code)
	rw.status = code
	rw.wroteHeader = true
}

func NewLogMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Error(
						"panic occured",
						slog.Any("error", err),
						slog.String("trace", string(debug.Stack())),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			log.Info(
				"request served",
				slog.Int("status", wrapped.Status()),
				slog.String("method", r.Method),
				slog.String("path", r.URL.EscapedPath()),
				slog.Duration("duration", time.Since(start)),
			)
		}

		return http.HandlerFunc(fn)
	}
}
