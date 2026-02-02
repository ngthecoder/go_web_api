package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewBadRequestError(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func NewNotFoundError(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func NewInternalServerError(message string, err error) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Err:        err,
	}
}

func NewConflictError(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}

func NewUnauthorizedError(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

func NewMethodNotAllowedError() *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusMethodNotAllowed,
		Message:    "Method not allowed",
	}
}

func WriteHTTPError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if httpErr, ok := err.(*HTTPError); ok {
		w.WriteHeader(httpErr.StatusCode)
		json.NewEncoder(w).Encode(map[string]string{
			"error": httpErr.Message,
		})
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server error",
		})
	}
}
