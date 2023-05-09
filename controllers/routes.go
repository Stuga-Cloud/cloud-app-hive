package controllers

import (
	applicationControllers "cloud-app-hive/controllers/applications"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) *gin.Engine {
	router.POST("/applications", applicationControllers.CreateAndDeployApplicationController)
	router.GET("/applications/:namespace/:name/metrics", applicationControllers.GetMetricsByApplicationNameAndNamespaceController)
	//router.GET("/applications", applicationControllers.GetApplicationController) TODO
	//router.GET("/applications/:namespace/:name/logs", applicationControllers.GetLogsByApplicationNameAndNamespaceController) TODO
	//router.GET("/applications/:namespace/:name", applicationControllers.GetApplicationByNameAndNamespaceController) TODO
	//router.GET("/applications/:user_id", applicationControllers.GetApplicationByUserIdController) TODO
	//router.PUT("/applications/:namespace/:name", applicationControllers.UpdateApplicationByNameAndNamespaceController) TODO
	//router.DELETE("/applications/:namespace/:name", applicationControllers.DeleteApplicationByNameAndNamespaceController) TODO

	return router
}
