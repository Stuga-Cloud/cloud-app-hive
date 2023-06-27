package controllers

import (
	"cloud-app-hive/use_cases"
	applicationsUseCases "cloud-app-hive/use_cases/applications"
	namespaceUseCases "cloud-app-hive/use_cases/namespaces"
	"net/http"

	"cloud-app-hive/controllers/applications"
	"cloud-app-hive/controllers/cluster"
	"cloud-app-hive/controllers/namespaces"

	"github.com/gin-gonic/gin"
)

func InitRoutes(
	router *gin.Engine,
	createNamespaceUseCase namespaceUseCases.CreateNamespaceUseCase,
	findNamespaceByIDUseCase namespaceUseCases.FindNamespaceByIDUseCase,
	findNamespacesUseCase namespaceUseCases.FindNamespacesUseCase,
	findApplicationsUseCase applicationsUseCases.FindApplicationsUseCase,
	findApplicationByIDUseCase applicationsUseCases.FindApplicationByIDUseCase,
	createApplicationUseCase applicationsUseCases.CreateApplicationUseCase,
	updateApplicationUseCase applicationsUseCases.UpdateApplicationUseCase,
	deleteApplicationUseCase applicationsUseCases.DeleteApplicationUseCase,
	deployApplicationUseCase applicationsUseCases.DeployApplicationUseCase,
	undeployApplicationUseCase applicationsUseCases.UndeployApplicationUseCase,
	getApplicationLogsUseCase applicationsUseCases.GetApplicationLogsUseCase,
	getApplicationMetricsUseCase applicationsUseCases.GetApplicationMetricsUseCase,
	getApplicationStatusUseCase applicationsUseCases.GetApplicationStatusUseCase,
	fillApplicationsStatusUseCase applicationsUseCases.FillApplicationStatusUseCase,
	createNamespaceMembershipUseCase namespaceUseCases.CreateNamespaceMembershipUseCase,
	removeNamespaceMembershipUseCase namespaceUseCases.RemoveNamespaceMembershipUseCase,
	deleteNamespaceByIDUseCase namespaceUseCases.DeleteNamespaceByIDUseCase,
	updateNamespaceByIDUseCase namespaceUseCases.UpdateNamespaceByIDUseCase,
	getClusterMetricsUseCase use_cases.GetClusterMetricsUseCase,
) *gin.Engine {
	api := router.Group("/api/v1")
	{
		api.GET("/health", HealthCheck)
		namespaces.InitNamespacesRoutes(
			api,
			createNamespaceUseCase,
			findNamespaceByIDUseCase,
			findNamespacesUseCase,
			createNamespaceMembershipUseCase,
			removeNamespaceMembershipUseCase,
			deleteNamespaceByIDUseCase,
			updateNamespaceByIDUseCase,
			fillApplicationsStatusUseCase,
		)
		applications.InitApplicationsRoutes(
			api,
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
			fillApplicationsStatusUseCase,
		)
		cluster.InitClusterRoutes(
			api,
			getClusterMetricsUseCase,
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
