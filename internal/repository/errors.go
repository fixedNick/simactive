package repository

import "errors"

var (
	ErrAlreadyExists = errors.New("sim already exists")
	ErrNotFound      = errors.New("sim not found")
)
