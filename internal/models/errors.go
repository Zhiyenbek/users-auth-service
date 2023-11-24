package models

import "errors"

var (
	ErrInvalidInput          = errors.New("INVALID_INPUT")
	ErrInternalServer        = errors.New("INTERNAL_SERVER_ERROR")
	ErrInvalidPasswordFormat = errors.New("WRONG_PASSWORD_FORMAT")
	ErrWrongPassword         = errors.New("WRONG_PASSWORD")
	ErrInvalidToken          = errors.New("INVALID_TOKEN")
	ErrTokenExpired          = errors.New("TOKEN_EXPIRED")
)
