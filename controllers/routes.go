package controllers

import (
	applicationsUseCases "cloud-app-hive/use_cases/applications"
	namespaceUseCases "cloud-app-hive/use_cases/namespaces"
	"net/http"

	"cloud-app-hive/controllers/applications"
	"cloud-app-hive/controllers/namespaces"

	"github.com/gin-gonic/gin"
)

func InitRoutes(
	router *gin.Engine,
	createNamespaceUseCase namespaceUseCases.CreateNamespaceUseCase,
	findNamespaceByIDUseCase namespaceUseCases.FindNamespaceByIDUseCase,
	findNamespacesUseCase namespaceUseCases.FindNamespacesUseCase,
	findApplicationsUseCase applicationsUseCases.FindApplicationsUseCase,
	createApplicationUseCase applicationsUseCases.CreateApplicationUseCase,
	updateApplicationUseCase applicationsUseCases.UpdateApplicationUseCase,
	deleteApplicationUseCase applicationsUseCases.DeleteApplicationUseCase,
	deployApplicationUseCase applicationsUseCases.DeployApplicationUseCase,
	undeployApplicationUseCase applicationsUseCases.UndeployApplicationUseCase,
	getApplicationLogsUseCase applicationsUseCases.GetApplicationLogsUseCase,
	getApplicationMetricsUseCase applicationsUseCases.GetApplicationMetricsUseCase,
	getApplicationStatusUseCase applicationsUseCases.GetApplicationStatusUseCase,
) *gin.Engine {
	api := router.Group("/api/v1")
	{
		api.GET("/health", HealthCheck)
		namespaces.InitNamespacesRoutes(api, createNamespaceUseCase, findNamespaceByIDUseCase, findNamespacesUseCase)
		applications.InitApplicationsRoutes(
			api,
			findApplicationsUseCase,
			createApplicationUseCase,
			updateApplicationUseCase,
			deleteApplicationUseCase,
			deployApplicationUseCase,
			undeployApplicationUseCase,
			getApplicationLogsUseCase,
			getApplicationMetricsUseCase,
			getApplicationStatusUseCase,
		)
	}

	return router
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description checks the health of the API
// @ID health-check
// @Tags Health
// @Produce  json
// @Success 200 {string} string	"pong"
// @Router /health [get].
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
