package server

import (
	"net/http"

	chttp "github.com/kaustubhbabar5/gh-api-client/pkg/http"
)

// registers routes for health check endpoints
func (app *app) RegisterHealthRoutes() {
	app.router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		chttp.WriteJSON(w, http.StatusOK, nil)
	})
}
