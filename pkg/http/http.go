package http

import (
	"encoding/json"
	"net/http"
)

// writes response with header Content-Type : application/json and provided http status code and response body.
func WriteJSON(w http.ResponseWriter, statusCode int, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}
