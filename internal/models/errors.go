package models

import (
	"errors"
)

var (
	ErrNoRecord        = errors.New("models: no matching found")
	ErrInvalidEmail    = errors.New("models: invalid email")
	ErrInvalidPassword = errors.New("models: invalid password")
	ErrDuplicateEmail  = errors.New("models: duplicate email")
)
