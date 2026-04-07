package repositories

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrEmailConflict = errors.New("email already exists")
	ErrConflict      = errors.New("resource conflict")
)
