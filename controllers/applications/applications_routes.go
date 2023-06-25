package applications

import (
	"cloud-app-hive/use_cases/applications"
	"github.com/gin-gonic/gin"
)

func InitApplicationsRoutes(
	router *gin.RouterGroup,
	findApplicationsUseCase applications.FindApplicationsUseCase,
	findApplicationByIDUseCase applications.FindApplicationByIDUseCase,
	createApplicationUseCase applications.CreateApplicationUseCase,
	updateApplicationUseCase applications.UpdateApplicationUseCase,
	deleteApplicationUseCase applications.DeleteApplicationUseCase,
	deployApplicationUseCase applications.DeployApplicationUseCase,
	undeployApplicationUseCase applications.UndeployApplicationUseCase,
	getApplicationLogsUseCase applications.GetApplicationLogsUseCase,
	getApplicationMetricsUseCase applications.GetApplicationMetricsUseCase,
	getApplicationStatusUseCase applications.GetApplicationStatusUseCase,
	fillApplicationStatusUseCase applications.FillApplicationStatusUseCase,
) {
	applicationController := NewApplicationController(
		findApplicationsUseCase,
		findApplicationByIDUseCase,
		createApplicationUseCase,
		updateApplicationUseCase,
		deleteApplicationUseCase,
		deployApplicationUseCase,
		undeployApplicationUseCase,
		getApplicationLogsUseCase,
		getApplicationMetricsUseCase,
		getApplicationStatusUseCase,
		fillApplicationStatusUseCase,
	)
	router.GET("/applications", applicationController.FindApplicationsController)
	router.POST("/applications", applicationController.CreateAndDeployApplicationController)
	router.GET("/applications/:id", applicationController.FindApplicationByIDController)
	router.PUT("/applications/:id", applicationController.UpdateApplicationByNameAndNamespaceController)
	router.GET("/applications/:id/metrics", applicationController.GetMetricsByApplicationNameAndNamespaceController)
	router.GET("/applications/:id/logs", applicationController.GetLogsByApplicationNameAndNamespaceController)
	router.GET("/applications/:id/status", applicationController.GetStatusByApplicationNameAndNamespaceController)
	router.DELETE("/applications/:id", applicationController.DeleteApplicationByIDController)
}
