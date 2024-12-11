package app

import (
	"context"
	"example/internal/config"
	"example/internal/http/v1"
	"example/internal/repository"
	"example/internal/service"
	"fmt"
	"net/http"
	"time"
)

type App struct {
	server *http.Server
}

func New(cfg *config.Config) (*App, error) {

	// logger

	// db

	mux := http.NewServeMux()

	repository := repository.NewRepository("inited db")
	service := service.NewService(repository)
	controllers := v1.NewControllers(service)

	// mux

	srv := &http.Server{
		Addr:              fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
		ReadHeaderTimeout: 800 * time.Millisecond,
		ReadTimeout:       800 * time.Millisecond,
		Handler:           mux,
	}

	return &App{server: srv}, nil
}

func (a *App) Run() error {

	// migrate

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(err.Error())
		}
	}()

	return nil
}

func (a *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Server gracefully stopped")

	//db

	return nil
}
