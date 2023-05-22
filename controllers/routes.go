package controllers

import (
	"cloud-app-hive/controllers/applications"
	"cloud-app-hive/controllers/namespaces"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) *gin.Engine {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	namespaces.InitNamespacesRoutes(router)
	applications.InitApplicationsRoutes(router)
	return router
}
