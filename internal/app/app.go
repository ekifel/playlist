package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ekifel/playlist/internal/api"
	"github.com/ekifel/playlist/internal/config"
	"github.com/ekifel/playlist/internal/postgres"
	"github.com/ekifel/playlist/internal/repository"
	"github.com/ekifel/playlist/internal/server"
	"github.com/ekifel/playlist/internal/service"
	"github.com/sirupsen/logrus"
)

func Run() {
	cfg, err := config.Init()
	if err != nil {
		logrus.Errorf("error initializing config: %v", err.Error())

		return
	}

	// Dependencies
	db, err := postgres.NewPostgresClient(&cfg.Postgres)
	if err != nil {
		logrus.Errorf("error initializing postgres client: %v", err.Error())

		return
	}

	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		Repos: repos,
	})
	handlers := api.NewHandler(services)

	srv := server.NewServer(cfg, handlers.Init())

	err = services.Playlist.Run()
	if err != nil {
		logrus.Errorf("error running playlist: %v", err.Error())

		return
	}

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logrus.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logrus.Errorf("failed to stop server: %v", err)
	}

	if err := db.Close(context.Background()); err != nil {
		logrus.Errorf("error occurred while closing connection to db: %s", err.Error())
	}
}
