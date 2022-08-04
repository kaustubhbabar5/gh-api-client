package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	cerrors "github.com/kaustubhbabar5/gh-api-client/pkg/errors"
	chttp "github.com/kaustubhbabar5/gh-api-client/pkg/http"
	"go.uber.org/zap"
)

type Handler struct {
	logger      *zap.Logger
	userService Service
}

func NewHandler(logger *zap.Logger, userService Service) *Handler {
	return &Handler{logger, userService}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		chttp.WriteJSON(
			w,
			http.StatusBadRequest,
			// TODO:
			map[string]any{"error": "failed to read request body, please refer API docs"},
		)
		return
	}

	var reqBodyMap map[string][]string

	err = json.Unmarshal(reqBody, &reqBodyMap)
	if err != nil {
		chttp.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]any{"error": "bad request, please refer API docs"},
		)
		return
	}

	//TODO: add validations for request body

	usernames := reqBodyMap["usernames"]

	users, notFoundUsers, errs := h.userService.GetUsers(usernames)
	// h.logger.Sugar().Info(users, errs, notFoundUsers)

	chttp.WriteJSON(w, http.StatusOK, map[string]any{
		"users":           users,
		"not_found_users": notFoundUsers,
		"errors":          cerrors.JSONErrs(errs),
	})
}
