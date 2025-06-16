package graphql

import (
	"errors"
	"time"

	"github.com/ductong169z/shorten-url/internal/models"
)

// Input types
type ShortenURLInput struct {
	OriginalURL string `json:"originalURL"`
	ShortCode   string `json:"shortCode,omitempty"`
}

// Response types
type ShortURLResponse struct {
	ID          uint64 `json:"id"`
	OriginalURL string `json:"originalURL"`
	ShortCode   string `json:"shortCode"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	FullShortURL string `json:"fullShortURL"`
}

// Validate validates the shorten URL input
func (s *ShortenURLInput) Validate() error {
	if s.OriginalURL == "" {
		return errors.New("original URL is required")
	}
	return nil
}

// Helper functions
func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Convert from model to response
func FromShortURLModel(shortURL *models.ShortURL, appDomain string) *ShortURLResponse {
	if shortURL == nil {
		return nil
	}

	fullShortURL := appDomain + "/" + shortURL.ShortCode

	return &ShortURLResponse{
		ID:          shortURL.ID,
		OriginalURL: shortURL.OriginalURL,
		ShortCode:   shortURL.ShortCode,
		CreatedAt:   formatTime(shortURL.CreatedAt),
		UpdatedAt:   formatTime(shortURL.UpdatedAt),
		FullShortURL: fullShortURL,
	}
}
