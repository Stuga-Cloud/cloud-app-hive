package controllers

import (
	"cloud-app-hive/controllers/applications/dto"
	"cloud-app-hive/domain"
	"cloud-app-hive/repositories"
	"cloud-app-hive/use_cases/applications"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// ApiError from https://github.com/go-playground/validator/issues/559
type ApiError struct {
	Param   string
	Message string
}

func CreateAndDeployApplicationController(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "Bearer "+os.Getenv("API_KEY") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var createApplicationDto dto.CreateApplicationDto
	if err := c.ShouldBindJSON(&createApplicationDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	err := dto.ValidateCreateApplicationDto(createApplicationDto)
	if err != nil {
		c.JSON(400, gin.H{
			"validation-errors": err.Error(),
		})
		return
	}

	deployApplicationUseCase := applications.DeployApplicationUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	deployApplication := domain.DeployApplication{
		Name:                      createApplicationDto.Name,
		Image:                     createApplicationDto.Image,
		Namespace:                 createApplicationDto.Namespace,
		Port:                      createApplicationDto.Port,
		ApplicationType:           createApplicationDto.ApplicationType,
		EnvironmentVariables:      nil,
		Secrets:                   nil,
		ContainerSpecifications:   domain.ApplicationContainerSpecifications{},
		ScalabilitySpecifications: domain.ApplicationScalabilitySpecifications{},
	}
	err = deployApplicationUseCase.Execute(deployApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while deploying": err.Error()})
		return
	}

	// TODO -> Create application in database
	createApplicationUseCase := applications.CreateApplicationUseCase{}
	application, err := createApplicationUseCase.Execute(deployApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while upserting in database": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("App %s deployed", application.Name),
	})
}

// GetMetricsByApplicationNameAndNamespaceController returns the metrics of an application by name and namespace in query params
func GetMetricsByApplicationNameAndNamespaceController(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "Bearer "+os.Getenv("API_KEY") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	applicationNamespace := c.Param("namespace")
	applicationName := c.Param("name")
	if applicationNamespace == "" || applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})
		return
	}

	// TODO -> Get metrics from kubernetes
	getMetricsByApplicationNameAndNamespaceUseCase := applications.GetApplicationMetricsUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	application := domain.GetApplicationMetrics{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	metrics, err := getMetricsByApplicationNameAndNamespaceUseCase.Execute(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while getting metrics": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}

// GetLogsByApplicationNameAndNamespaceController returns the logs of an application by name and namespace in query params
func GetLogsByApplicationNameAndNamespaceController(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "Bearer "+os.Getenv("API_KEY") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	applicationNamespace := c.Param("namespace")
	applicationName := c.Param("name")
	if applicationNamespace == "" || applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})
		return
	}

	// TODO -> Get logs from kubernetes
	getLogsByApplicationNameAndNamespaceUseCase := applications.GetApplicationLogsUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	application := domain.GetApplicationLogs{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	logs, err := getLogsByApplicationNameAndNamespaceUseCase.Execute(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while getting logs": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}
