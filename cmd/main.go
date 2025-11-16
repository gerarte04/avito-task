package main

import (
	"avito-task/internal/config"
	"avito-task/internal/usecases/service"
	pkgConfig "avito-task/pkg/config"
	"avito-task/pkg/database/postgres"
	"avito-task/pkg/shutdown"
	"context"
	"errors"
	"log"
	"time"

	httpapp "avito-task/internal/app/http"
	repo "avito-task/internal/repository/postgres"

	"golang.org/x/sync/errgroup"
)

func main() {
	appFlags := pkgConfig.ParseFlags()
	var cfg config.Config
	pkgConfig.MustLoadConfig(appFlags.ConfigPath, &cfg)

	log.Printf("[INFO] Service is starting")

	pool, err := postgres.NewPostgresPool(cfg.PostgresCfg)
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect PostgreSQL: %s", err.Error())
	}

	defer pool.Close()
	log.Printf("[INFO] Connected to PostgreSQL successfully")

	teamRepo := repo.NewTeamRepo(pool)
	userRepo := repo.NewUserRepo(pool)
	prRepo := repo.NewPullRequestRepo(pool)

	teamSvc := service.NewTeamService(pool, teamRepo, userRepo)
	userSvc := service.NewUserService(userRepo, prRepo)
	prSvc := service.NewPullRequestService(pool, prRepo, userRepo)

	httpApp := httpapp.New(
		cfg.HTTPCfg,
		cfg.PathCfg,
		cfg.SvcCfg,
		teamSvc,
		userSvc,
		prSvc,
	)

	log.Printf("[INFO] All services were created successfully")

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return shutdown.ListenSignal(ctx)
	})

	g.Go(func() error {
		return httpApp.Run()
	})

	g.Go(func() error {
		<-ctx.Done()
		log.Printf("[INFO] Shutdown signal received, stopping server")

		const ctxTimeExceed = 3 * time.Second

		shutdownCtx, cancel := context.WithTimeout(context.Background(), ctxTimeExceed)
		defer cancel()
		return httpApp.Stop(shutdownCtx)
	})

	err = g.Wait()
	if err != nil && !errors.Is(err, shutdown.ErrOSSignal) {
		log.Printf("[INFO] Exit reason: %s", err.Error())
	}
}
