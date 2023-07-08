package domain

import (
	"github.com/google/uuid"
	"time"
)

type ShortenedURL struct {
	ID          uuid.UUID `json:"id"`
	LongURL     string    `json:"long_url"`
	ShortURL    string    `json:"short_url"`
	QRCodeURL   string    `json:"qr_code_url"`
	CustomAlias string    `json:"custom_alias,omitempty"`
	Clicks      int       `json:"clicks"`
	LastAccess  string    `json:"last_access"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
