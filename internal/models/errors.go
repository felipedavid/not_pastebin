package models

import "errors"

var (
	ErrNoRecord       = errors.New("models: no record found")
	ErrDuplicateEmail = errors.New("models: email already in use")
)
