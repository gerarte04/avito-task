package server

import (
	"avito-task/pkg/config"
	"net/http"
)

func CreateServer(handler http.Handler, cfg config.HTTPConfig) error {
    s := &http.Server{
        Handler: handler,
        Addr: cfg.Address,
        ReadTimeout: cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
        IdleTimeout: cfg.IdleTimeout,
    }

    return s.ListenAndServe()
}
