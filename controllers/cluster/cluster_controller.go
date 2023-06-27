package cluster

import (
	"cloud-app-hive/use_cases"
	"net/http"

	controllerValidators "cloud-app-hive/controllers/validators"

	"github.com/gin-gonic/gin"
)

type ClusterController struct {
	getClusterMetricsUseCase use_cases.GetClusterMetricsUseCase
}

func NewClusterController(
	getClusterMetricsUseCase use_cases.GetClusterMetricsUseCase,
) ClusterController {
	return ClusterController{
		getClusterMetricsUseCase: getClusterMetricsUseCase,
	}
}

func (clusterController ClusterController) GetClusterMetricsController(c *gin.Context) {
	if !controllerValidators.ValidateAuthorizationToken(c) {
		controllerValidators.Unauthorized(c)
		return
	}

	clusterMetrics, err := clusterController.getClusterMetricsUseCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metrics": clusterMetrics})
}
