package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/kaustubhbabar5/gh-api-client/adapters/cache"
	"github.com/kaustubhbabar5/gh-api-client/pkg/config"
	"github.com/kaustubhbabar5/gh-api-client/pkg/github"

	"go.uber.org/zap"
)

type application struct {
	logger       *zap.Logger
	router       *mux.Router
	cache        *redis.Client
	githubClient github.Client
}

// constructs http server.
func New(config *config.Config, logger *zap.Logger) (*http.Server, error) {
	cache, err := cache.NewRedisClient(config.RedisURL, config.RedisPasswordKey)
	if err != nil {
		return nil, fmt.Errorf("cache.NewRedisCache :%w", err)
	}

	router := mux.NewRouter()

	httpClient := &http.Client{
		Timeout: time.Duration(config.ReadTimeout) * time.Second,
	}

	githubClient := github.NewClient(httpClient, config.GithubAuthTokenKey)

	app := application{
		logger,
		router,
		cache,
		githubClient,
	}

	app.RegisterHealthRoutes()
	app.RegisterUserRoutes()

	return &http.Server{
		Addr:              config.Host + ":" + config.Port,
		ReadTimeout:       time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(config.WriteTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(config.ReadHeaderTimeout) * time.Second,
		Handler:           app.router,
	}, nil
}
