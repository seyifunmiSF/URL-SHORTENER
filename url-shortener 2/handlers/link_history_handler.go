package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"url-shortener/repositories"
)

func LinkHistoryHandler(repository *repositories.ShortenedURLRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve link history from the database
		urlDataList, err := repository.GetAllURLData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Transform the URL data into a response format
		var linkHistory []string
		for _, urlData := range urlDataList {
			linkHistory = append(linkHistory, urlData.ShortURL)
		}

		c.JSON(http.StatusOK, gin.H{
			"linkHistory": linkHistory,
			"data":        urlDataList,
		})
	}
}
