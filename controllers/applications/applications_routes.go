package applications

import (
	"cloud-app-hive/use_cases/applications"
	"github.com/gin-gonic/gin"
)

func InitApplicationsRoutes(
	router *gin.RouterGroup,
	createApplicationUseCase applications.CreateApplicationUseCase,
	updateApplicationUseCase applications.UpdateApplicationUseCase,
	deleteApplicationUseCase applications.DeleteApplicationUseCase,
	deployApplicationUseCase applications.DeployApplicationUseCase,
	undeployApplicationUseCase applications.UndeployApplicationUseCase,
	getApplicationLogsUseCase applications.GetApplicationLogsUseCase,
	getApplicationMetricsUseCase applications.GetApplicationMetricsUseCase,
	getApplicationStatusUseCase applications.GetApplicationStatusUseCase,
) {
	applicationController := NewApplicationController(
		createApplicationUseCase,
		updateApplicationUseCase,
		deleteApplicationUseCase,
		deployApplicationUseCase,
		undeployApplicationUseCase,
		getApplicationLogsUseCase,
		getApplicationMetricsUseCase,
		getApplicationStatusUseCase,
	)
	router.POST("/applications", applicationController.CreateAndDeployApplicationController)
	router.PUT("/applications/:id", applicationController.UpdateApplicationByNameAndNamespaceController)
	router.GET("/applications/:namespace/:name/metrics", applicationController.GetMetricsByApplicationNameAndNamespaceController)
	router.GET("/applications/:namespace/:name/logs", applicationController.GetLogsByApplicationNameAndNamespaceController)
	router.GET("/applications/:namespace/:name/status", applicationController.GetStatusByApplicationNameAndNamespaceController)
	// router.GET("/applications", applicationControllers.GetApplicationController) TODO when database is implemented
	// router.GET("/applications/:namespace/:name", applicationControllers.GetApplicationByNameAndNamespaceController) TODO
	// router.GET("/applications/:user_id", applicationControllers.GetApplicationByUserIdController) TODO
	// router.PUT("/applications/:namespace/:name", applicationControllers.UpdateApplicationByNameAndNamespaceController) TODO
	router.DELETE("/applications/:id", applicationController.DeleteApplicationByNameAndNamespaceController)
}
