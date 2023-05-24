package validators

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ValidateAuthorizationToken(c *gin.Context) bool {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "Bearer "+os.Getenv("API_KEY") {
		return false
	}
	return true
}

func Unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}
