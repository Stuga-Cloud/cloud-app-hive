package controllers

import (
	"cloud-app-hive/controllers/applications/dto"
	"cloud-app-hive/domain"
	"cloud-app-hive/repositories"
	"cloud-app-hive/use_cases/applications"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"net/http"
	"os"
)

// ApiError from https://github.com/go-playground/validator/issues/559
type ApiError struct {
	Param   string
	Message string
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return fmt.Sprintf("Minimum length is %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", fe.Param())
	case "email":
		return "Invalid email address"
	default:
		return "Invalid input - " + fe.Tag()
	}
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

	// TODO -> Deploy application
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
	output, err := deployApplicationUseCase.Execute(deployApplication)
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
		"message":  fmt.Sprintf("App %s deployed", application.Name),
		"output":   output,
		"cpuUsage": 0,
		"ramUsage": 0,
		"time":     0,
	})
}

// GetMetricsByApplicationNameAndNamespaceController returns the metrics of an application by name and namespace in query params
func GetMetricsByApplicationNameAndNamespaceController(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader != "Bearer "+os.Getenv("API_KEY") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	appNamespace := c.Param("namespace")
	appName := c.Param("name")
	if appNamespace == "" || appName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})
		return
	}
	println("APP NAME: " + appName)
	println("APP NAMESPACE: " + appNamespace)

	// TODO -> Get metrics from kubernetes
	getMetricsByApplicationNameAndNamespaceUseCase := applications.GetApplicationMetricsUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	metrics, err := getMetricsByApplicationNameAndNamespaceUseCase.Execute(appName, appNamespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while getting metrics": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}
