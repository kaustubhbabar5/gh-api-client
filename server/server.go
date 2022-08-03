package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kaustubhbabar5/gh-api-client/adapters/cache"
	"github.com/kaustubhbabar5/gh-api-client/pkg/config"
	"github.com/kaustubhbabar5/gh-api-client/pkg/github"
	"go.uber.org/zap"
)

type app struct {
	router       *mux.Router
	cache        cache.Cache
	githubClient github.Client
}

//constructs http server
func New(config *config.Config, logger *zap.Logger) (*http.Server, error) {
	cache, err := cache.NewRedisCache(config.RedisUrl, config.RedisPassword)
	if err != nil {
		return nil, fmt.Errorf("cache.NewRedisCache :%w", err)
	}

	router := mux.NewRouter()

	httpClient := &http.Client{
		Timeout: time.Duration(config.ReadTimeout) * time.Second,
	}

	githubClient := github.NewClient(httpClient, config.GithubAuthToken)

	app := app{
		router,
		cache,
		githubClient,
	}

	app.RegisterHealthRoutes()

	return &http.Server{
		Addr:         config.Host + ":" + config.Port,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
		Handler:      app.router,
	}, nil
}
