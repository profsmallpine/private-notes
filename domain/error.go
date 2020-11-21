package domain

import "errors"

var (
	ErrNotValid     = errors.New("invalid type")
	ErrNoUser       = errors.New("no user found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrUnexpected   = errors.New("unexpected error")
)
