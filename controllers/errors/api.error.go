package errors

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Date        string      `json:"date"`
	StatusCode  int         `json:"status_code"`
	Name        string      `json:"name"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
	Context     interface{} `json:"context"`
}

func NewApiError(statusCode int, name string, message string, description string, c *gin.Context, context interface{}) ApiError {
	apiContext := map[string]interface{}{
		"method":      c.Request.Method,
		"host":        c.Request.Host,
		"path":        c.Request.URL.Path,
		"body":        c.Request.Body,
		"params":      c.Params,
		"headers":     c.Request.Header,
		"query":       c.Request.URL.Query(),
		"contentType": c.ContentType(),
		"clientIP":    c.ClientIP(),
		"fullPath":    c.FullPath(),
	}
	// append the apiContext to the context
	context.(map[string]interface{})["apiContext"] = apiContext
	err := ApiError{
		Date:        time.Now().Format(time.RFC3339),
		StatusCode:  statusCode,
		Name:        name,
		Message:     message,
		Description: description,
		Context:     context,
	}
	fmt.Println(err.String())
	return err
}

func (e ApiError) String() string {
	return fmt.Sprintf("%s - %d - %s - %s - %s - %v", e.Date, e.StatusCode, e.Name, e.Message, e.Description, e.Context)
}
