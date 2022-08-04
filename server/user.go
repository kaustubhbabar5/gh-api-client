package server

import "github.com/kaustubhbabar5/gh-api-client/internal/user"

func (app *application) RegisterUserRoutes() {
	userCache := user.NewCache(app.logger, app.cache)
	userService := user.NewService(app.logger, userCache, app.githubClient)
	userHandler := user.NewHandler(app.logger, app.validator, userService)

	app.router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
}
