package validators

import (
	"github.com/gin-gonic/gin"
)

func ValidateBodyIsNotNullNorEmpty(c *gin.Context) bool {
	if c.Request.Body == nil || c.Request.ContentLength == 0 {
		return false
	}
	return true
}

func BodyIsNullOrEmptyResponse(c *gin.Context) {
	c.JSON(400, gin.H{"error": "Request body must not be empty"})
}
