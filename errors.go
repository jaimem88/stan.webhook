package webhook

import (
	"fmt"
	"net/http"
)

var (
	errBadRequest                 = &Error{Code: http.StatusBadRequest, ErrorMessage: "Could not decode request: JSON parsing failed"}
	errBadRequestMalformedPayload = &Error{Code: http.StatusBadRequest, ErrorMessage: "Could not decode request: Malformed Payload"}
	errInternalServerError        = &Error{Code: http.StatusInternalServerError, ErrorMessage: "Something went wrong :("}
)

// Error describes custom error that can be used for logging and to write the response inside the handler
type Error struct {
	message      string // msg used for logging purposes
	Code         int    `json:"code,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

// Error implements error interface
func (e *Error) Error() string {
	return fmt.Sprintf("code: %d message: %s", e.Code, e.ErrorMessage)
}

func (e *Error) msg(m string) *Error {
	e.message = m
	return e
}
