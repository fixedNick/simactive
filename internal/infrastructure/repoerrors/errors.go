package repoerrors

import "errors"

var (
	ErrAlreadyExists = errors.New("Already exists")
	ErrNotFound      = errors.New("Not found")
)
