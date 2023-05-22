package applications

import (
	"cloud-app-hive/controllers/applications/requests"
	"cloud-app-hive/controllers/applications/responses"
	"cloud-app-hive/controllers/errors"
	"cloud-app-hive/controllers/validators"
	"cloud-app-hive/database"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/repositories"
	"cloud-app-hive/use_cases/applications"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateAndDeployApplicationController godoc
// @Summary Creates in database and deploys an application
// @Description creates in database and deploys an application on the cloud
// @ID create-and-deploy-application
// @Tags Applications
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization Token"
// @Param createApplicationRequest body requests.CreateApplicationRequest true "Create Application Request"
// @Success 200 {object} responses.CreateApplicationResponse
// @Failure 400 {object} errors.ApiError
// @Router /applications [post]
func CreateAndDeployApplicationController(c *gin.Context) {
	if validators.ValidateAuthorizationToken(c) == false {
		validators.Unauthorized(c)
		return
	}

	var createApplicationRequest requests.CreateApplicationRequest
	if err := c.ShouldBindJSON(&createApplicationRequest); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewApiError(
			http.StatusBadRequest,
			"validation_errors",
			fmt.Errorf("error while binding json: %w", err).Error(),
			"Check if the request body is correct (see the swagger documentation)",
			c,
			map[string]interface{}{
				"body": c.Request.Body,
			},
		))
		return
	}

	err := requests.ValidateCreateApplicationRequest(createApplicationRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"validation-errors": fmt.Errorf("error while validating create application request: %w", err).Error(),
		})
		return
	}

	db, err := database.ConnectToDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	createApplicationUseCase := applications.CreateApplicationUseCase{
		ApplicationRepository: repositories.GORMApplicationRepository{
			Database: db,
		},
		NamespaceRepository: repositories.GORMNamespaceRepository{
			Database: db,
		},
	}
	createApplication := commands.CreateApplication{
		Name:                      createApplicationRequest.Name,
		Image:                     createApplicationRequest.Image,
		NamespaceID:               createApplicationRequest.NamespaceID,
		UserID:                    createApplicationRequest.UserID,
		Port:                      createApplicationRequest.Port,
		ApplicationType:           createApplicationRequest.ApplicationType,
		EnvironmentVariables:      createApplicationRequest.EnvironmentVariables,
		Secrets:                   createApplicationRequest.Secrets,
		ContainerSpecifications:   createApplicationRequest.ContainerSpecifications,
		ScalabilitySpecifications: createApplicationRequest.ScalabilitySpecifications,
	}
	application, namespace, err := createApplicationUseCase.Execute(createApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deployApplicationUseCase := applications.DeployApplicationUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	applyApplication := commands.ApplyApplication{
		Name:                      application.Name,
		Image:                     application.Image,
		Namespace:                 namespace.Name,
		Port:                      application.Port,
		ApplicationType:           application.ApplicationType,
		EnvironmentVariables:      *application.EnvironmentVariables,
		Secrets:                   *application.Secrets,
		ContainerSpecifications:   *application.ContainerSpecifications,
		ScalabilitySpecifications: *application.ScalabilitySpecifications,
	}
	err = deployApplicationUseCase.Execute(applyApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.CreateApplicationResponse{
		Message:     fmt.Sprintf("App %s deployed", application.Name),
		Application: responses.ApplicationDomainToResponse(application),
	})
}

func UpdateApplicationByNameAndNamespaceController(c *gin.Context) {
	if validators.ValidateAuthorizationToken(c) == false {
		validators.Unauthorized(c)
		return
	}

	var updateApplicationRequest requests.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&updateApplicationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	applicationID := c.Param("id")
	userID := c.Query("user_id")
	if applicationID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": "Application ID url param and user_id query param are required"})
		return
	}

	err := requests.ValidateUpdateApplicationRequest(updateApplicationRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"validation-errors": err.Error(),
		})
		return
	}

	db, err := database.ConnectToDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	updateApplicationUseCase := applications.UpdateApplicationUseCase{
		ApplicationRepository: repositories.GORMApplicationRepository{
			Database: db,
		},
		NamespaceRepository: repositories.GORMNamespaceRepository{
			Database: db,
		},
	}
	updateApplication := commands.UpdateApplication{
		UserID:                    userID,
		Description:               updateApplicationRequest.Description,
		Image:                     updateApplicationRequest.Image,
		Port:                      updateApplicationRequest.Port,
		ApplicationType:           updateApplicationRequest.ApplicationType,
		EnvironmentVariables:      nil,
		Secrets:                   nil,
		ContainerSpecifications:   updateApplicationRequest.ContainerSpecifications,
		ScalabilitySpecifications: updateApplicationRequest.ScalabilitySpecifications,
	}
	application, namespace, err := updateApplicationUseCase.Execute(applicationID, updateApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"context": "While updating application in database",
			"body":    updateApplicationRequest,
		})
		return
	}

	deployApplicationUseCase := applications.DeployApplicationUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	applyApplication := commands.ApplyApplication{
		Name:                      application.Name,
		Image:                     application.Image,
		Namespace:                 namespace.Name,
		Port:                      application.Port,
		ApplicationType:           application.ApplicationType,
		EnvironmentVariables:      *application.EnvironmentVariables,
		Secrets:                   *application.Secrets,
		ContainerSpecifications:   *application.ContainerSpecifications,
		ScalabilitySpecifications: *application.ScalabilitySpecifications,
	}
	err = deployApplicationUseCase.Execute(applyApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     fmt.Sprintf("App %s deployed", application.Name),
		"application": application,
	})
}

// DeleteApplicationByNameAndNamespaceController deletes an application by name and namespace in query params
func DeleteApplicationByNameAndNamespaceController(c *gin.Context) {
	if validators.ValidateAuthorizationToken(c) == false {
		validators.Unauthorized(c)
		return
	}

	applicationID := c.Param("id")
	userID := c.Query("user_id")
	if applicationID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID url param and 'user_id' query param must be provided"})
		return
	}

	db, err := database.ConnectToDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	deleteApplicationUseCase := applications.DeleteApplicationUseCase{
		ApplicationRepository: repositories.GORMApplicationRepository{
			Database: db,
		},
	}
	deleteApplication := commands.DeleteApplication{
		ID:     applicationID,
		UserID: userID,
	}
	deletedApplication, err := deleteApplicationUseCase.Execute(deleteApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	undeployApplicationUseCase := applications.UndeployApplicationUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	unapplyApplication := commands.UnapplyApplication{
		Name:      deletedApplication.Name,
		Namespace: deletedApplication.Namespace.Name,
	}
	err = undeployApplicationUseCase.Execute(unapplyApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     fmt.Sprintf("App %s deleted in namespace %s", deletedApplication.Name, deletedApplication.Namespace.Name),
		"application": deletedApplication,
	})
}

// GetMetricsByApplicationNameAndNamespaceController returns the metrics of an application by name and namespace in query params
func GetMetricsByApplicationNameAndNamespaceController(c *gin.Context) {
	if validators.ValidateAuthorizationToken(c) == false {
		validators.Unauthorized(c)
		return
	}

	applicationNamespace := c.Param("namespace")
	applicationName := c.Param("name")
	if applicationNamespace == "" || applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})
		return
	}

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
	if validators.ValidateAuthorizationToken(c) == false {
		validators.Unauthorized(c)
		return
	}

	applicationNamespace := c.Param("namespace")
	applicationName := c.Param("name")
	if applicationNamespace == "" || applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})
		return
	}

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

// GetStatusByApplicationNameAndNamespaceController returns the status of an application by name and namespace in query params
func GetStatusByApplicationNameAndNamespaceController(c *gin.Context) {
	if validators.ValidateAuthorizationToken(c) == false {
		validators.Unauthorized(c)
		return
	}

	applicationNamespace := c.Param("namespace")
	applicationName := c.Param("name")
	if applicationNamespace == "" || applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})
		return
	}

	// TODO -> Get status from kubernetes
	getStatusByApplicationNameAndNamespaceUseCase := applications.GetApplicationStatusUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	application := commands.GetApplicationStatus{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	status, err := getStatusByApplicationNameAndNamespaceUseCase.Execute(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
