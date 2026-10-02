package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
)


type App struct {
	router http.Handler
	rdb *redis.Client


}

func New() *App {
	app := &App{
		router: loadRoutes(),
		rdb:   redis.NewClient(&redis.Options{}),
	}

	return  app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr: ":3000",
		Handler: a.router,
	}

    err := a.rdb.Ping().Err()

	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)

	}

	fmt.Println("Starting service")


	err = server.ListenAndServe()

	if err != nil {
		return  fmt.Errorf("failed to start the server: %w", err)
	}

	return  nil

}

