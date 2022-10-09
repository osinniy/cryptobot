package api

import (
	"fmt"
	"net/http"
)

type apiError struct {
	Code    int
	Message string
}

var (
	errUserNotFound  = newError("user not found", http.StatusNotFound)
	errMissingUserId = newError("missing user_id", http.StatusBadRequest)
)

func newError(text string, statusCode int) apiError {
	return apiError{
		Code:    statusCode,
		Message: fmt.Sprintf("%d: %s", statusCode, text),
	}
}

func (s Server) e(w http.ResponseWriter, apiErr apiError) {
	w.WriteHeader(apiErr.Code)
	_, err := w.Write([]byte(apiErr.Message))
	s.handleWrErr(err, &w)
}

func (s Server) handleWrErr(err error, w *http.ResponseWriter) {
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to write response")
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}
}
