package controllers

import (
	"cloud-app-hive/controllers/applications/dto"
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
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

	createApplicationUseCase := applications.CreateApplicationUseCase{}
	createApplication := commands.CreateApplication{
		Name:                      createApplicationDto.Name,
		Image:                     createApplicationDto.Image,
		Namespace:                 createApplicationDto.Namespace,
		Port:                      createApplicationDto.Port,
		ApplicationType:           createApplicationDto.ApplicationType,
		EnvironmentVariables:      createApplicationDto.EnvironmentVariables,      // TODO ADAPTERS
		Secrets:                   createApplicationDto.Secrets,                   // TODO ADAPTERS
		ContainerSpecifications:   createApplicationDto.ContainerSpecifications,   // TODO ADAPTERS
		ScalabilitySpecifications: createApplicationDto.ScalabilitySpecifications, // TODO ADAPTERS
	}
	application, err := createApplicationUseCase.Execute(createApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deployApplicationUseCase := applications.DeployApplicationUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	applyApplication := commands.ApplyApplication{
		Name:                      createApplicationDto.Name,
		Image:                     createApplicationDto.Image,
		Namespace:                 createApplicationDto.Namespace,
		Port:                      createApplicationDto.Port,
		ApplicationType:           createApplicationDto.ApplicationType,
		EnvironmentVariables:      createApplicationDto.EnvironmentVariables,      // TODO ADAPTERS
		Secrets:                   createApplicationDto.Secrets,                   // TODO ADAPTERS
		ContainerSpecifications:   createApplicationDto.ContainerSpecifications,   // TODO ADAPTERS
		ScalabilitySpecifications: createApplicationDto.ScalabilitySpecifications, // TODO ADAPTERS
	}
	err = deployApplicationUseCase.Execute(applyApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("App %s deployed", application.Name),
	})
}

func UpdateApplicationByNameAndNamespaceController(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "Bearer "+os.Getenv("API_KEY") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var updateApplicationDto dto.UpdateApplicationDto
	if err := c.ShouldBindJSON(&updateApplicationDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	err := dto.ValidateUpdateApplicationDto(updateApplicationDto)
	if err != nil {
		c.JSON(400, gin.H{
			"validation-errors": err.Error(),
		})
		return
	}

	deployApplicationUseCase := applications.DeployApplicationUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	applyApplication := commands.ApplyApplication{
		Name:                      updateApplicationDto.Name,
		Image:                     updateApplicationDto.Image,
		Namespace:                 updateApplicationDto.Namespace,
		Port:                      updateApplicationDto.Port,
		ApplicationType:           updateApplicationDto.ApplicationType,
		EnvironmentVariables:      nil,                                           // TODO ADAPTERS
		Secrets:                   nil,                                           // TODO ADAPTERS
		ContainerSpecifications:   domain.ApplicationContainerSpecifications{},   // TODO ADAPTERS
		ScalabilitySpecifications: domain.ApplicationScalabilitySpecifications{}, // TODO ADAPTERS
	}
	err = deployApplicationUseCase.Execute(applyApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updateApplicationUseCase := applications.UpdateApplicationUseCase{}
	updateApplication := commands.UpdateApplication{
		Name:                      updateApplicationDto.Name,
		Image:                     updateApplicationDto.Image,
		Namespace:                 updateApplicationDto.Namespace,
		Port:                      updateApplicationDto.Port,
		ApplicationType:           updateApplicationDto.ApplicationType,
		EnvironmentVariables:      nil,                                           // TODO ADAPTERS
		Secrets:                   nil,                                           // TODO ADAPTERS
		ContainerSpecifications:   domain.ApplicationContainerSpecifications{},   // TODO ADAPTERS
		ScalabilitySpecifications: domain.ApplicationScalabilitySpecifications{}, // TODO ADAPTERS
	}
	application, err := updateApplicationUseCase.Execute(updateApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "context": "While updating application in database", "data": updateApplicationDto})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("App %s updated", application.Name),
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
	application := commands.GetApplicationMetrics{
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
	application := commands.GetApplicationLogs{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	logs, err := getLogsByApplicationNameAndNamespaceUseCase.Execute(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}

// DeleteApplicationByNameAndNamespaceController deletes an application by name and namespace in query params
func DeleteApplicationByNameAndNamespaceController(c *gin.Context) {
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

	undeployApplicationUseCase := applications.UndeployApplicationUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	unapplyApplication := commands.UnapplyApplication{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	err := undeployApplicationUseCase.Execute(unapplyApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deleteApplicationUseCase := applications.DeleteApplicationUseCase{}
	deleteApplication := commands.DeleteApplication{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	_, err = deleteApplicationUseCase.Execute(deleteApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("App %s deleted", applicationName),
	})
}
