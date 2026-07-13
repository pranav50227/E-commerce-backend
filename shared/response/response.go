package response

import (
	"github.com/gin-gonic/gin"
)

// JSON sends a standard successful JSON response.
func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// Error sends a structured JSON error response.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
