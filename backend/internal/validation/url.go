package validation

import (
	"net/url"

	"qrgen/internal/apperrors"
)

func ValidateHTTPURL(rawURL string) error {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return apperrors.ErrInvalidURL
	}

	if parsed.Host == "" {
		return apperrors.ErrInvalidURL
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return apperrors.ErrInvalidURL
	}

	return nil
}
