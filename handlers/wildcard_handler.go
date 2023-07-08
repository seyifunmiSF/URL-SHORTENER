package handlers

import (
	"github.com/gin-gonic/gin"
)

func WildcardHandler(c *gin.Context) {
	filepath := c.Param("filepath")

	// Specify the directory path where your image files are located
	projectDirectory := "./images/"

	// Create the complete file path by concatenating the project directory and the requested filepath
	fullPath := projectDirectory + filepath

	// Serve the file using the complete file path
	c.File(fullPath)
}
