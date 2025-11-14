package handlers

import (
	pkgMiddleware "avito-task/pkg/http/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	SwaggerPath = "/swagger/*"
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

func WithSwagger() RouterOption {
	return func(r chi.Router) {
		r.Get(SwaggerPath, httpSwagger.WrapHandler)
	}
}
