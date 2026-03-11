package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Render(w http.ResponseWriter, err error) {
	code := toHTTPStatus(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(HTTPError{
		Code:    code,
		Message: err.Error(),
	})
}

func toHTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound // 404
	case errors.Is(err, ErrConflict):
		return http.StatusConflict // 409
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden // 403
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest // 400
	case errors.Is(err, ErrUnauth):
		return http.StatusUnauthorized // 401
	default:
		return http.StatusInternalServerError // 500
	}
}
