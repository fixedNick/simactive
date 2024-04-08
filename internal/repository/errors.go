package repository

import "errors"

var (
	ErrSimAlreadyExists = errors.New("sim already exists")
	ErrSimNotFound      = errors.New("sim not found")
	ErrServiceNotFound  = errors.New("service not found")
)
