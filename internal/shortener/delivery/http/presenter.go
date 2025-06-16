package http

import (
	"regexp"
	"strings"

	"github.com/ductong169z/shorten-url/internal/models"
	"github.com/ductong169z/shorten-url/internal/shortener"
)

type ShortURLResponse struct {
	ID          uint64  `json:"id"`
	OriginalURL string  `json:"original_url"`
	ShortURL    string  `json:"short_url"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	ExpiredAt   *string `json:"expired_at,omitempty"`
	ClickCount  uint    `json:"click_count"`
}

func FromShortURLModel(url *models.ShortURL, domain string) ShortURLResponse {
	var expiredAt *string
	if url.ExpiredAt != nil {
		v := url.ExpiredAt.Format("2006-01-02 15:04:05")
		expiredAt = &v
	}
	return ShortURLResponse{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		ShortURL:    domain + "/" + url.ShortCode,
		CreatedAt:   url.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   url.UpdatedAt.Format("2006-01-02 15:04:05"),
		ExpiredAt:   expiredAt,
		ClickCount:  url.ClickCount,
	}
}

type ShortenRequest struct {
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code,omitempty"`
}

// Validate checks the OriginalURL prefix
func (r *ShortenRequest) Validate() error {
	if !(len(r.OriginalURL) > 0 && (strings.HasPrefix(r.OriginalURL, "http://") || strings.HasPrefix(r.OriginalURL, "https://"))) {
		return shortener.ErrInvalidOriginalURL
	}
	if r.ShortCode != "" {
		if len(r.ShortCode) < 4 || len(r.ShortCode) > 16 || !regexp.MustCompile(`^[a-zA-Z0-9_-]{4,16}$`).MatchString(r.ShortCode) {
			return shortener.ErrInvalidShortCode
		}
	}

	return nil
}

type ShortenResponse struct {
	ID          uint64  `json:"id"`
	OriginalURL string  `json:"original_url"`
	ShortCode   string  `json:"short_code"`
	ShortURL    string  `json:"short_url"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	ExpiredAt   *string `json:"expired_at,omitempty"`
	ClickCount  uint    `json:"click_count"`
}
