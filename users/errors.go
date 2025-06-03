package users

import "errors"

var (
	ErrInvalidParamFormat = errors.New("Invalid parameter format")
	ErrProfileExists      = errors.New("Profile already exists, update existing one")
	ErrInvalidDateFormat  = errors.New("Invalid date format, expected YYYY-MM-DD")
)
