// Package shortener provides core error definitions and utilities for the URL shortener domain.
// It defines domain-specific error variables and error-to-HTTP status mapping for consistent error handling.
package shortener

import (
	"errors"
	"net/http"
)

const (
	// shortCodeAlreadyExists is returned when a short code already exists in the system.
	shortCodeAlreadyExists = "short code already exists"
	// shortCodeNotFound is returned when a requested short code does not exist.
	shortCodeNotFound = "short code not found"
	// shortCodeExpired is returned when a requested short code has expired.
	shortCodeExpired = "short code expired"
	// invalidOriginalURL is returned when the original URL is invalid.
	invalidOriginalURL = "invalid original URL"
	// invalidShortCode is returned when the provided short code is invalid.
	invalidShortCode = "invalid short code"
)

var (
	// ErrShortCodeAlreadyExists indicates that the short code already exists.
	ErrShortCodeAlreadyExists = errors.New(shortCodeAlreadyExists)
	// ErrShortCodeNotFound indicates that the short code was not found.
	ErrShortCodeNotFound = errors.New(shortCodeNotFound)
	// ErrShortCodeExpired indicates that the short code has expired.
	ErrShortCodeExpired = errors.New(shortCodeExpired)
	// ErrInvalidOriginalURL indicates that the original URL is invalid.
	ErrInvalidOriginalURL = errors.New(invalidOriginalURL)
	// ErrInvalidShortCode indicates that the provided short code is invalid.
	ErrInvalidShortCode = errors.New(invalidShortCode)
)

// MapError maps a domain error to an HTTP status code and message.
// It provides a unified way to translate domain errors to HTTP responses.
func MapError(err error) (status int, message string) {
	switch {
	case errors.Is(err, ErrShortCodeAlreadyExists):
		return http.StatusConflict, shortCodeAlreadyExists
	case errors.Is(err, ErrShortCodeNotFound):
		return http.StatusNotFound, shortCodeNotFound
	case errors.Is(err, ErrShortCodeExpired):
		return http.StatusGone, shortCodeExpired
	case errors.Is(err, ErrInvalidOriginalURL):
		return http.StatusBadRequest, invalidOriginalURL
	case errors.Is(err, ErrInvalidShortCode):
		return http.StatusBadRequest, invalidShortCode
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
