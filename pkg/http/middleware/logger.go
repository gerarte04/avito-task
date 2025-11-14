package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		next.ServeHTTP(ww, r)

		log.Printf(
			"[HTTP/1.1] | %s | %d | %s | %s | %s",
			r.RemoteAddr, ww.Status(), r.Method, r.URL.Path, time.Since(start).String(),
		)
	})
}
