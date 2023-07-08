package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"url-shortener/domain"
	"url-shortener/services"
)

func ShortenURLHandler(service *services.ShortenURLService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var urlData domain.ShortenedURL

		err := c.BindJSON(&urlData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Shorten the URL
		shortenedURL, err := service.ShortenURL(&urlData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":        shortenedURL,
			"qr_code_url": shortenedURL.QRCodeURL,
		})
	}
}

func RedirectHandler(service *services.ShortenURLService) gin.HandlerFunc {
	return func(c *gin.Context) {
		alias := c.Param("alias")
		// Perform rate limiting
		rateLimitDuration := 50
		maxRequestsPerDuration := 15
		// Check if rate limit is exceeded
		if service.PerformRateLimiting(alias, time.Duration(rateLimitDuration), int64(maxRequestsPerDuration)) {
			handleRateLimitError(c)
			return
		}

		url, err := service.RedirectOriginalURL(alias)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Invalid URL"})
		}
		c.Redirect(http.StatusMovedPermanently, url)
	}
}

func StatisticsHandler(service *services.ShortenURLService) gin.HandlerFunc {
	return func(c *gin.Context) {
		alias := c.Param("alias")
		statistics, err := service.GetURLStatistics(alias)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"shortenedURL": statistics.ShortURL,
			"originalURL":  statistics.LongURL,
			"visitCount":   statistics.Clicks,
			"qrCodeURL":    statistics.QRCodeURL,
		})
	}
}

func handleRateLimitError(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
}
