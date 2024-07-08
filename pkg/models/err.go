package models

import (
	"errors"
	"net/http"
)

var (
	//InternalServerError
	ErrInternalServer = errors.New("internal server error. please try again later")

	//NotFound
	ErrAccountNotFound = errors.New("account not found")

	//BadRequest
	ErrAccountExists = errors.New("account already exists")
	ErrBadRequest    = errors.New("bad request")

	//Unauthorized
	ErrTokenExpired       = errors.New("token expired")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var customErrors = map[error]int{
	ErrInternalServer: http.StatusInternalServerError,

	ErrAccountNotFound: http.StatusNotFound,

	ErrAccountExists: http.StatusBadRequest,
	ErrBadRequest:    http.StatusBadRequest,

	ErrTokenExpired:       http.StatusUnauthorized,
	ErrUnauthorized:       http.StatusUnauthorized,
	ErrInvalidCredentials: http.StatusUnauthorized,
}

func IsErrorCustom(err error) bool {
	for k := range customErrors {
		if errors.Is(err, k) {
			return true
		}
	}
	return false
}

func ErrorStatusCode(err error) int {
	for k, v := range customErrors {
		if errors.Is(err, k) {
			return v
		}
	}
	return http.StatusInternalServerError
}
