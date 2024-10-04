package db

import "errors"

// List of errors
var (
	ErrUnsupportedType = errors.New("infra: unsupported type")
	ErrNotFound        = errors.New("document: not found")
	ErrDuplicateKey    = errors.New("infra: duplicate key")
	ErrInvalidData     = errors.New("infra: invalid data")
)
