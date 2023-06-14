package applications

import (
	validators "cloud-app-hive/validators"
	"fmt"
	"net/http"

	"cloud-app-hive/controllers/applications/requests"
	"cloud-app-hive/controllers/applications/responses"
	"cloud-app-hive/controllers/errors"
	controllerValidators "cloud-app-hive/controllers/validators"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/repositories"
	"cloud-app-hive/use_cases/applications"

	"github.com/gin-gonic/gin"
)

type ApplicationController struct {
	findApplicationsUseCase      applications.FindApplicationsUseCase
	createApplicationUseCase     applications.CreateApplicationUseCase
	updateApplicationUseCase     applications.UpdateApplicationUseCase
	deleteApplicationUseCase     applications.DeleteApplicationUseCase
	deployApplicationUseCase     applications.DeployApplicationUseCase
	undeployApplicationUseCase   applications.UndeployApplicationUseCase
	getApplicationLogsUseCase    applications.GetApplicationLogsUseCase
	getApplicationMetricsUseCase applications.GetApplicationMetricsUseCase
	getApplicationStatusUseCase  applications.GetApplicationStatusUseCase
}

func NewApplicationController(
	findApplicationsUseCase applications.FindApplicationsUseCase,
	createApplicationUseCase applications.CreateApplicationUseCase,
	updateApplicationUseCase applications.UpdateApplicationUseCase,
	deleteApplicationUseCase applications.DeleteApplicationUseCase,
	deployApplicationUseCase applications.DeployApplicationUseCase,
	undeployApplicationUseCase applications.UndeployApplicationUseCase,
	getApplicationLogsUseCase applications.GetApplicationLogsUseCase,
	getApplicationMetricsUseCase applications.GetApplicationMetricsUseCase,
	getApplicationStatusUseCase applications.GetApplicationStatusUseCase,
) ApplicationController {
	return ApplicationController{
		findApplicationsUseCase:      findApplicationsUseCase,
		createApplicationUseCase:     createApplicationUseCase,
		updateApplicationUseCase:     updateApplicationUseCase,
		deleteApplicationUseCase:     deleteApplicationUseCase,
		deployApplicationUseCase:     deployApplicationUseCase,
		undeployApplicationUseCase:   undeployApplicationUseCase,
		getApplicationLogsUseCase:    getApplicationLogsUseCase,
		getApplicationMetricsUseCase: getApplicationMetricsUseCase,
		getApplicationStatusUseCase:  getApplicationStatusUseCase,
	}
}

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
func (applicationController ApplicationController) CreateAndDeployApplicationController(c *gin.Context) {
	if !controllerValidators.ValidateAuthorizationToken(c) {
		controllerValidators.Unauthorized(c)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"validation-errors": fmt.Errorf("error while validating create application request: %w", err).Error(),
		})

		return
	}

	err = validators.ValidateEmail(createApplicationRequest.AdministratorEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"validation-errors": fmt.Errorf("error while validating create application request: %w", err).Error(),
		})
		return
	}
	createApplication := commands.CreateApplication{
		Name:                      createApplicationRequest.Name,
		Description:               createApplicationRequest.Description,
		Image:                     createApplicationRequest.Image,
		NamespaceID:               createApplicationRequest.NamespaceID,
		UserID:                    createApplicationRequest.UserID,
		Port:                      createApplicationRequest.Port,
		Zone:                      createApplicationRequest.Zone,
		ApplicationType:           createApplicationRequest.ApplicationType,
		EnvironmentVariables:      createApplicationRequest.EnvironmentVariables,
		Secrets:                   createApplicationRequest.Secrets,
		ContainerSpecifications:   createApplicationRequest.ContainerSpecifications,
		ScalabilitySpecifications: createApplicationRequest.ScalabilitySpecifications,
		AdministratorEmail:        createApplicationRequest.AdministratorEmail,
	}

	application, namespace, err := applicationController.createApplicationUseCase.Execute(createApplication)
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

// FindApplicationsController godoc
// @Summary Finds all applications
// @Description finds all applications
// @ID find-applications
// @Tags Applications
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization Token"
// @Success 200 {array} Application
// @Failure 400 {object} errorsApiError
// @Router /applications [get]
func (applicationController ApplicationController) FindApplicationsController(c *gin.Context) {
	if !controllerValidators.ValidateAuthorizationToken(c) {
		controllerValidators.Unauthorized(c)
		return
	}

	var findApplicationsRequest requests.FindApplicationsRequest
	if err := c.ShouldBindQuery(&findApplicationsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	findApplicationsCommand := commands.FindApplications{
		Name:            findApplicationsRequest.Name,
		Image:           findApplicationsRequest.Image,
		NamespaceID:     findApplicationsRequest.NamespaceID,
		ApplicationType: findApplicationsRequest.ApplicationType,
		IsAutoScaled:    findApplicationsRequest.IsAutoScaled,
		Page:            findApplicationsRequest.Page,
		Limit:           findApplicationsRequest.Limit,
	}
	foundApplications, err := applicationController.findApplicationsUseCase.Execute(findApplicationsCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"applications": foundApplications})
}

func (applicationController ApplicationController) UpdateApplicationByNameAndNamespaceController(c *gin.Context) {
	if !controllerValidators.ValidateAuthorizationToken(c) {
		controllerValidators.Unauthorized(c)
		return
	}

	var updateApplicationRequest requests.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&updateApplicationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	applicationID := c.Param("id")
	userID := c.Query("userId")
	if applicationID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": "Application ID url param and userId query param are required"})
		return
	}

	err := requests.ValidateUpdateApplicationRequest(updateApplicationRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"validation-errors": err.Error(),
		})
		return
	}

	err = validators.ValidateEmail(updateApplicationRequest.AdministratorEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"validation-errors": fmt.Errorf("error while validating update application request: %w", err).Error(),
		})
		return
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
		AdministratorEmail:        updateApplicationRequest.AdministratorEmail,
	}
	application, namespace, err := applicationController.updateApplicationUseCase.Execute(applicationID, updateApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"context": "While updating application in database",
			"body":    updateApplicationRequest,
		})
		return
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
	err = applicationController.deployApplicationUseCase.Execute(applyApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     fmt.Sprintf("App %s deployed", application.Name),
		"application": application,
	})
}

// DeleteApplicationByIDController deletes an application by name and namespace in query params
func (applicationController ApplicationController) DeleteApplicationByIDController(c *gin.Context) {
	if !controllerValidators.ValidateAuthorizationToken(c) {
		controllerValidators.Unauthorized(c)
		return
	}

	applicationID := c.Param("id")
	userID := c.Query("userId")
	if applicationID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID url param and 'userId' query param must be provided"})
		return
	}

	deleteApplication := commands.DeleteApplication{
		ID:     applicationID,
		UserID: userID,
	}
	deletedApplication, err := applicationController.deleteApplicationUseCase.Execute(deleteApplication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unapplyApplication := commands.UnapplyApplication{
		Name:      deletedApplication.Name,
		Namespace: deletedApplication.Namespace.Name,
	}
	err = applicationController.undeployApplicationUseCase.Execute(unapplyApplication)
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
func (applicationController ApplicationController) GetMetricsByApplicationNameAndNamespaceController(c *gin.Context) {
	if !controllerValidators.ValidateAuthorizationToken(c) {
		controllerValidators.Unauthorized(c)
		return
	}

	applicationNamespace := c.Param("namespace")
	applicationName := c.Param("name")
	if applicationNamespace == "" || applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})
		return
	}

	application := commands.GetApplicationMetrics{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	metrics, err := applicationController.getApplicationMetricsUseCase.Execute(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while getting metrics": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}

// GetLogsByApplicationNameAndNamespaceController returns the logs of an application by name and namespace in query params
func (applicationController ApplicationController) GetLogsByApplicationNameAndNamespaceController(c *gin.Context) {
	if !controllerValidators.ValidateAuthorizationToken(c) {
		controllerValidators.Unauthorized(c)
		return
	}

	applicationNamespace := c.Param("namespace")
	applicationName := c.Param("name")
	if applicationNamespace == "" || applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})

		return
	}

	application := commands.GetApplicationLogs{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	logs, err := applicationController.getApplicationLogsUseCase.Execute(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}

// GetStatusByApplicationNameAndNamespaceController returns the status of an application by name and namespace in query params
func (applicationController ApplicationController) GetStatusByApplicationNameAndNamespaceController(c *gin.Context) {
	if !controllerValidators.ValidateAuthorizationToken(c) {
		controllerValidators.Unauthorized(c)
		return
	}

	applicationNamespace := c.Param("namespace")
	applicationName := c.Param("name")
	if applicationNamespace == "" || applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})
		return
	}

	// TODO -> Get status from kubernetes
	application := commands.GetApplicationStatus{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	status, err := applicationController.getApplicationStatusUseCase.Execute(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
