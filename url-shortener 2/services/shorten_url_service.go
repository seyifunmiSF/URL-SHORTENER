package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"time"
	"url-shortener/domain"
	"url-shortener/repositories"
	"url-shortener/utils"
)

type ShortenURLService struct {
	shortenerUrlRepository repositories.ShortenedURLRepository
	redisRepository        repositories.RedisRepository
	config                 *domain.Config
	qrCodeService          domain.QRCodeService
}

func NewShortenURLService(repository repositories.ShortenedURLRepository, redisRepository repositories.RedisRepository, config *domain.Config, qrCodeService domain.QRCodeService) *ShortenURLService {
	return &ShortenURLService{
		shortenerUrlRepository: repository,
		redisRepository:        redisRepository,
		config:                 config,
		qrCodeService:          qrCodeService,
	}
}

func (s *ShortenURLService) ShortenURL(urlData *domain.ShortenedURL) (*domain.ShortenedURL, error) {
	// Validate the long URL
	if err := validateURL(urlData.LongURL); err != nil {
		return nil, err
	}

	// Generate a unique short URL
	shortURL, uniqueID, err := s.generateUniqueShortURL()
	if err != nil {
		return nil, err
	}

	urlData.ID = uuid.New()

	urlData.Clicks = 0
	urlData.CreatedAt = time.Now()
	urlData.UpdatedAt = time.Now()

	// Check if a custom alias is provided
	if urlData.CustomAlias != "" {
		// Check if the custom alias is available
		available, err := s.checkCustomAliasAvailability(urlData.CustomAlias)
		if err != nil {
			return nil, err
		}
		if !available {
			return nil, fmt.Errorf("custom alias '%s' is already taken", urlData.CustomAlias)
		}
		// Set the custom alias in the URL data
		urlData.ShortURL = s.config.BaseURL + "/s/" + urlData.CustomAlias
	} else {
		// Set the generated short URL in the URL data
		urlData.ShortURL = shortURL
		urlData.CustomAlias = uniqueID
	}

	// Generate the QR code for the shortened URL
	qrCodeURL, err := s.qrCodeService.GenerateQRCode(urlData.ShortURL)
	if err != nil {
		return nil, err
	}

	// Set the QR code URL in the URL data
	urlData.QRCodeURL = qrCodeURL

	// Save the URL data to the database
	err = s.shortenerUrlRepository.SaveURLData(urlData)
	if err != nil {
		return nil, err
	}

	// Save the URL data to Redis
	err = s.redisRepository.SaveURLData(urlData)
	if err != nil {
		return nil, err
	}

	return urlData, nil
}

func (s *ShortenURLService) RedirectOriginalURL(alias string) (string, error) {
	urlData, err := s.redisRepository.GetURLDataByAlias(alias)
	if err != nil {
		// Cache miss, fetch from the database
		urlData, err = s.shortenerUrlRepository.GetURLDataByAlias(alias)
		if err != nil {
			return "", errors.New("URL not found")
		}
		// Cache the fetched data for future use
		if err := s.redisRepository.SaveURLData(urlData); err != nil {
			// Log or handle the error as appropriate
		}
	}
	err = s.UpdateAnalytics(alias)
	if err != nil {
		return "", err
	}
	if urlData == nil {
		return "", errors.New("URL not found")
	}
	if &urlData.LongURL == nil || urlData.LongURL == "" {
		return "", errors.New("URL not found")
	}
	return urlData.LongURL, nil
}

// GetURLStatistics Get statistics for a single URL
func (s *ShortenURLService) GetURLStatistics(alias string) (*domain.ShortenedURL, error) {
	urlData, err := s.shortenerUrlRepository.GetURLDataByAlias(alias)
	if err != nil {
		return nil, errors.New("URL not found")
	}
	return urlData, nil
}

func (s *ShortenURLService) checkCustomAliasAvailability(alias string) (bool, error) {
	// Query the database to check if the alias exists
	exists, err := s.shortenerUrlRepository.CheckCustomAliasExists(alias)
	if err != nil {
		return false, err
	}
	// If the alias exists, it is not available
	if exists {
		return false, nil
	}
	return true, nil
}

func (s *ShortenURLService) generateUniqueShortURL() (string, string, error) {
	// Generate a short URL using a unique identifier
	uniqueID, err := utils.GenerateUniqueID()
	if err != nil {
		return "", "", err
	}

	shortURL := s.config.BaseURL + "/s/" + uniqueID
	return shortURL, uniqueID, nil
}

func (s *ShortenURLService) UpdateAnalytics(alias string) error {
	urlData, err := s.shortenerUrlRepository.GetURLDataByAlias(alias)
	if err != nil {
		return err
	}

	// Update analytics
	urlData.Clicks++
	urlData.LastAccess = time.Now().Format(time.RFC3339)

	// Save the updated URL data
	err = s.shortenerUrlRepository.UpdateURLData(urlData)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShortenURLService) PerformRateLimiting(alias string, duration time.Duration, maxRequests int64) bool {
	return s.redisRepository.PerformRateLimiting(alias, duration, maxRequests)
}

func validateURL(longURL string) error {
	// Parse the URL
	_, err := url.ParseRequestURI(longURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %s", longURL)
	}
	return nil
}
