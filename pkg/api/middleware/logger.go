package middleware

import (
	"net/http"
	"time"

	"github.com/pressly/chi/middleware"
	"go.uber.org/zap"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
// Credit: github.com/treastech/logger
func Logger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				l.Info("Served",
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.String("latency", time.Since(t1).String()),
					zap.Int("status", ww.Status()))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
