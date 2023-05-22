package controllers

import (
	"cloud-app-hive/controllers/applications"
	"cloud-app-hive/controllers/namespaces"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) *gin.Engine {
	router.GET("/health", HealthCheck)

	namespaces.InitNamespacesRoutes(router)
	applications.InitApplicationsRoutes(router)
	return router
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description checks the health of the API
// @ID health-check
// @Tags Health
// @Produce  json
// @Success 200 {string} string	"pong"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
