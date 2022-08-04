package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
	cerrors "github.com/kaustubhbabar5/gh-api-client/pkg/errors"
	chttp "github.com/kaustubhbabar5/gh-api-client/pkg/http"
	"go.uber.org/zap"
)

type Handler struct {
	logger    *zap.Logger
	validator *validator.Validate

	userService Service
}

func NewHandler(logger *zap.Logger, validator *validator.Validate, userService Service) *Handler {
	return &Handler{logger, validator, userService}
}

// swagger:route GET /users  GetUsers
// Get Users information
//
// responses:
//  500: InternalServerError
//  200: GetUsers
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

	var req GetUsersRequest

	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		chttp.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]any{"error": "bad request, please refer API docs"},
		)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		chttp.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]any{"error": "bad request, please refer API docs"},
		)
		return
	}
	//TODO: add validations for request

	usernames := req.Usernames

	users, notFoundUsers, errs := h.userService.GetUsers(usernames)
	// h.logger.Sugar().Info(users, errs, notFoundUsers)
	if len(users) == 0 {
		h.logger.Warn("failed to fetch users", zap.Errors("errors", errs))
		chttp.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	chttp.WriteJSON(w, http.StatusOK, GetUsersResponse{
		Users:         users,
		NotFoundUsers: notFoundUsers,
		Errors:        cerrors.JSONErrs(errs),
	})
}
