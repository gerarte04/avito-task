package handlers

import (
	pkgMiddleware "avito-task/pkg/http/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterOption func(r chi.Router)

func RouteHandlers(r chi.Router, apiPath string, opts ...RouterOption) {
	r.Route(apiPath, func(r chi.Router) {
		for _, opt := range opts {
			opt(r)
		}
	})
}

func WithLogger() RouterOption {
	return func(r chi.Router) {
		r.Use(pkgMiddleware.Logger)
	}
}

func WithRecovery() RouterOption {
	return func(r chi.Router) {
		r.Use(middleware.Recoverer)
	}
}

func WithSwagger(path string, fsRoot string) RouterOption {
	srv := http.FileServer(http.Dir(fsRoot))

	return func(r chi.Router) {
		r.Handle(path + "/*", http.StripPrefix(path + "/", srv))
	}
}

func WithHealthHandler() RouterOption {
	return func(r chi.Router) {
		r.Mount("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("OK"))
		}))
	}
}
