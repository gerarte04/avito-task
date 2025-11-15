package http

import (
	apihttp "avito-task/internal/api/http"
	"avito-task/internal/config"
	"avito-task/internal/usecases"
	"avito-task/pkg/http/handlers"
	"context"
	"log"
	"net/http"

	pkgConfig "avito-task/pkg/config"

	"github.com/go-chi/chi/v5"
)

type App struct {
	server *http.Server
}

func New(
	httpCfg pkgConfig.HTTPConfig,
	pathCfg config.PathConfig,
	teamSvc usecases.TeamService,
	userSvc usecases.UserService,
	prSvc usecases.PullRequestService,
) *App {
	teamHandler := apihttp.NewTeamHandler(teamSvc, pathCfg)
	userHandler := apihttp.NewUserHandler(userSvc, pathCfg)
	prHandler := apihttp.NewPullRequestHandler(prSvc, pathCfg)

	router := chi.NewRouter()
	handlers.RouteHandlers(router, pathCfg.APIPath,
		handlers.WithLogger(),
		handlers.WithRecovery(),
		handlers.WithSwagger("/swagger"),
		teamHandler.WithTeamHandlers(),
		userHandler.WithUserHandlers(),
		prHandler.WithPRHandlers(),
	)

	srv := &http.Server{
		Handler:      router,
		Addr:         httpCfg.Address,
		ReadTimeout:  httpCfg.ReadTimeout,
		WriteTimeout: httpCfg.WriteTimeout,
		IdleTimeout:  httpCfg.IdleTimeout,
	}

	return &App{
		server: srv,
	}
}

func (a *App) Run() error {
	const op = "http.App.Run"

	log.Printf("[INFO] %s: starting server at %s...", op, a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	const op = "http.App.Stop"

	log.Printf("[INFO] %s: http server shutting down", op)
	return a.server.Shutdown(ctx)
}
