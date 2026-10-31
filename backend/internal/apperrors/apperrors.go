package apperrors

import "errors"

var (
	ErrNotFound   = errors.New("qr code not found")
	ErrInvalidURL = errors.New("invalid URL: only http and https URLs are supported")
)
