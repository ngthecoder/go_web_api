package errors

import (
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

func WriteHTTPError(w http.ResponseWriter, err error) {
	if httpErr, ok := err.(*HTTPError); ok {
		http.Error(w, httpErr.Message, httpErr.StatusCode)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
