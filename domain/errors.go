package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrRecordNotFound      = errors.New("record not found")
	ErrBadRequest          = errors.New("bad request")
)

type ErrorResponse struct {
	Message string `json:"message"`
}
