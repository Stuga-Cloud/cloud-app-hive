package controllers

import (
	lambdaControllers "cloud-app-hive/controllers/applications"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) *gin.Engine {
	router.POST("/create-application", lambdaControllers.CreateContainerController)

	return router
}
