package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type WildcardService struct {
	// Add any required dependencies or repositories here
}

func NewWildcardService() *WildcardService {
	// Initialize any dependencies or repositories here
	return &WildcardService{}
}

func (s *WildcardService) WildcardHandler(c *gin.Context) {
	filepath := c.Param("filepath")
	// Add your logic to handle the wildcard route here
	// Example: Serve the file with the given filepath
	http.ServeFile(c.Writer, c.Request, filepath)
}
