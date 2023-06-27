package cluster

import (
	"cloud-app-hive/use_cases"

	"github.com/gin-gonic/gin"
)

func InitClusterRoutes(
	router *gin.RouterGroup,
	getClusterMetricsUseCase use_cases.GetClusterMetricsUseCase,
) {
	clusterController := NewClusterController(
		getClusterMetricsUseCase,
	)
	router.GET("/cluster/metrics", clusterController.GetClusterMetricsController)
}
