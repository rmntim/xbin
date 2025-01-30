package errors

import "errors"

var (
	ErrNotFound = errors.New("bin not found")
	ErrExpired  = errors.New("bin expired")
)
