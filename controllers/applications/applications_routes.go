package applications

import (
	"github.com/gin-gonic/gin"
)

func InitApplicationsRoutes(router *gin.RouterGroup) {
	router.POST("/applications", CreateAndDeployApplicationController)
	router.PUT("/applications/:id", UpdateApplicationByNameAndNamespaceController)
	router.GET("/applications/:namespace/:name/metrics", GetMetricsByApplicationNameAndNamespaceController)
	router.GET("/applications/:namespace/:name/logs", GetLogsByApplicationNameAndNamespaceController)
	router.GET("/applications/:namespace/:name/status", GetStatusByApplicationNameAndNamespaceController)
	// router.GET("/applications", applicationControllers.GetApplicationController) TODO when database is implemented
	// router.GET("/applications/:namespace/:name", applicationControllers.GetApplicationByNameAndNamespaceController) TODO
	// router.GET("/applications/:user_id", applicationControllers.GetApplicationByUserIdController) TODO
	// router.PUT("/applications/:namespace/:name", applicationControllers.UpdateApplicationByNameAndNamespaceController) TODO
	router.DELETE("/applications/:id", DeleteApplicationByNameAndNamespaceController)
}
