package applications

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/use_cases"
	validators "cloud-app-hive/validators"
	"fmt"
	"net/http"
	"os"
	"strconv"

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
	findApplicationsUseCase       applications.FindApplicationsUseCase
	findApplicationByIDUseCase    applications.FindApplicationByIDUseCase
	createApplicationUseCase      applications.CreateApplicationUseCase
	updateApplicationUseCase      applications.UpdateApplicationUseCase
	deleteApplicationUseCase      applications.DeleteApplicationUseCase
	deployApplicationUseCase      applications.DeployApplicationUseCase
	undeployApplicationUseCase    applications.UndeployApplicationUseCase
	getApplicationLogsUseCase     applications.GetApplicationLogsUseCase
	getApplicationMetricsUseCase  applications.GetApplicationMetricsUseCase
	getApplicationStatusUseCase   applications.GetApplicationStatusUseCase
	fillApplicationsStatusUseCase applications.FillApplicationStatusUseCase
	getClusterMetricsUseCase      use_cases.GetClusterMetricsUseCase
}

func NewApplicationController(
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
	fillApplicationsStatusUseCase applications.FillApplicationStatusUseCase,
	getClusterMetricsUseCase use_cases.GetClusterMetricsUseCase,
) ApplicationController {
	return ApplicationController{
		findApplicationsUseCase:       findApplicationsUseCase,
		findApplicationByIDUseCase:    findApplicationByIDUseCase,
		createApplicationUseCase:      createApplicationUseCase,
		updateApplicationUseCase:      updateApplicationUseCase,
		deleteApplicationUseCase:      deleteApplicationUseCase,
		deployApplicationUseCase:      deployApplicationUseCase,
		undeployApplicationUseCase:    undeployApplicationUseCase,
		getApplicationLogsUseCase:     getApplicationLogsUseCase,
		getApplicationMetricsUseCase:  getApplicationMetricsUseCase,
		getApplicationStatusUseCase:   getApplicationStatusUseCase,
		fillApplicationsStatusUseCase: fillApplicationsStatusUseCase,
		getClusterMetricsUseCase:      getClusterMetricsUseCase,
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
		fmt.Println("Error while binding json when creating application: ", err)
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
		fmt.Println("Error while validating create application request: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"validation-errors": fmt.Errorf("error while validating create application request: %w", err).Error(),
		})

		return
	}

	err = validators.ValidateEmail(createApplicationRequest.AdministratorEmail)
	if err != nil {
		fmt.Println("Error while validating create application request: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"validation-errors": fmt.Errorf("error while validating create application request: %w", err).Error(),
		})
		return
	}

	// Check cluster state before creating application
	clusterState, err := applicationController.getClusterMetricsUseCase.Execute()
	if err != nil {
		fmt.Println("Error while getting cluster state: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if clusterState == nil {
		fmt.Println("Cluster state is nil")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cluster state is nil"})
		return
	}

	stopDeployingApplicationWhenClusterNodesUsageIsAbovePercentageStr := os.Getenv("STOP_DEPLOYING_APPLICATION_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE")
	if stopDeployingApplicationWhenClusterNodesUsageIsAbovePercentageStr == "" {
		fmt.Println("STOP_DEPLOYING_APPLICATION_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE is not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "STOP_DEPLOYING_APPLICATION_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE is not set"})
		return
	}
	stopDeployingApplicationWhenPercentageOfNodesExceededUsageStr := os.Getenv("STOP_DEPLOYING_APPLICATION_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE")
	if stopDeployingApplicationWhenPercentageOfNodesExceededUsageStr == "" {
		fmt.Println("STOP_DEPLOYING_APPLICATION_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE is not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "STOP_DEPLOYING_APPLICATION_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE is not set"})
	}

	stopDeployingApplicationWhenClusterNodesUsageIsAbovePercentage, err := strconv.ParseFloat(stopDeployingApplicationWhenClusterNodesUsageIsAbovePercentageStr, 64)
	if err != nil {
		fmt.Println("Error when convert STOP_DEPLOYING_APPLICATION_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE to float64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when convert STOP_DEPLOYING_APPLICATION_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE to float64"})
		return

	}
	stopDeployingApplicationWhenPercentageOfNodesExceededUsage, err := strconv.ParseFloat(stopDeployingApplicationWhenPercentageOfNodesExceededUsageStr, 64)
	if err != nil {
		fmt.Println("Error when convert STOP_DEPLOYING_APPLICATION_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE to float64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when convert STOP_DEPLOYING_APPLICATION_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE to float64"})
	}

	// Check if cluster is not exceeding its limits
	if domain.DoesPartOfNodesExceedCPUOrMemoryUsage(
		stopDeployingApplicationWhenClusterNodesUsageIsAbovePercentage,
		stopDeployingApplicationWhenPercentageOfNodesExceededUsage,
		clusterState.NodesComputedUsages,
	) {
		fmt.Println("Cluster is exceeding its limits")
		c.JSON(http.StatusInsufficientStorage, gin.H{"error": "Cluster is exceeding its limits"})
		return
	}

	createApplication := commands.CreateApplication{
		Name:                      createApplicationRequest.Name,
		Description:               createApplicationRequest.Description,
		Image:                     createApplicationRequest.Image,
		Registry:                  createApplicationRequest.Registry,
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
		fmt.Println("Error while creating application: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deployApplicationUseCase := applications.DeployApplicationUseCase{
		ContainerManagerRepository: repositories.KubernetesContainerManagerRepository{},
	}
	applyApplication := commands.ApplyApplication{
		Name:                      application.Name,
		Image:                     application.Image,
		Registry:                  application.Registry,
		Namespace:                 namespace.Name,
		Port:                      application.Port,
		ApplicationType:           application.ApplicationType,
		EnvironmentVariables:      *application.EnvironmentVariables,
		Secrets:                   *application.Secrets,
		ContainerSpecifications:   *application.ContainerSpecifications,
		ScalabilitySpecifications: application.ScalabilitySpecifications.Data(),
	}
	err = deployApplicationUseCase.Execute(applyApplication)
	if err != nil {
		fmt.Println("Error while deploying application: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.CreateApplicationResponse{
		Message:     fmt.Sprintf("App %s deployed", application.Name),
		Application: *application,
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

func (applicationController ApplicationController) FindApplicationByIDController(c *gin.Context) {
	if !controllerValidators.ValidateAuthorizationToken(c) {
		controllerValidators.Unauthorized(c)
		return
	}

	applicationID := c.Param("id")
	if applicationID == "" {
		fmt.Println("Application ID url param is required")
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": "Application ID url param is required"})
		return
	}

	queryByUserID := c.Query("userId")
	if queryByUserID == "" {
		fmt.Println("userId query param is required")
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": "userId query param is required"})
		return
	}

	foundApplication, err := applicationController.findApplicationByIDUseCase.Execute(commands.FindApplicationByID{
		ApplicationID: applicationID,
		QueryByUserID: queryByUserID,
	})
	if err != nil {
		fmt.Println("Error while finding application by ID: ", err)
		if _, ok := err.(*errors.UnauthorizedToAccessNamespaceError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	applicationsWithStatus, err := applicationController.fillApplicationsStatusUseCase.Execute(foundApplication.Namespace.Name, []domain.Application{*foundApplication})
	if err != nil {
		fmt.Println("Error while filling applications status: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	foundApplication = &applicationsWithStatus[0]

	c.JSON(http.StatusOK, gin.H{"application": foundApplication})
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
		fmt.Println("Application ID url param and userId query param are required")
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": "Application ID url param and userId query param are required"})
		return
	}

	err := requests.ValidateUpdateApplicationRequest(updateApplicationRequest)
	if err != nil {
		fmt.Println("Error while validating update application request: ", err)
		c.JSON(400, gin.H{
			"validation-errors": err.Error(),
		})
		return
	}

	err = validators.ValidateEmail(updateApplicationRequest.AdministratorEmail)
	if err != nil {
		fmt.Println("Error while validating update application request: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"validation-errors": fmt.Errorf("error while validating update application request: %w", err).Error(),
		})
		return
	}
	updateApplication := commands.UpdateApplication{
		UserID:                    userID,
		Description:               updateApplicationRequest.Description,
		Image:                     updateApplicationRequest.Image,
		Registry:                  updateApplicationRequest.Registry,
		Port:                      updateApplicationRequest.Port,
		ApplicationType:           updateApplicationRequest.ApplicationType,
		EnvironmentVariables:      updateApplicationRequest.EnvironmentVariables,
		Secrets:                   updateApplicationRequest.Secrets,
		ContainerSpecifications:   updateApplicationRequest.ContainerSpecifications,
		ScalabilitySpecifications: updateApplicationRequest.ScalabilitySpecifications,
		AdministratorEmail:        updateApplicationRequest.AdministratorEmail,
	}
	application, namespace, err := applicationController.updateApplicationUseCase.Execute(applicationID, updateApplication, userID)
	if err != nil {
		if _, ok := err.(*errors.UnauthorizedToAccessNamespaceError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Error while updating application: ", err)
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
		Registry:                  application.Registry,
		Namespace:                 namespace.Name,
		Port:                      application.Port,
		ApplicationType:           application.ApplicationType,
		EnvironmentVariables:      *application.EnvironmentVariables,
		Secrets:                   *application.Secrets,
		ContainerSpecifications:   *application.ContainerSpecifications,
		ScalabilitySpecifications: application.ScalabilitySpecifications.Data(),
	}
	err = applicationController.deployApplicationUseCase.Execute(applyApplication)
	if err != nil {
		fmt.Println("Error while deploying application: ", err)
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
		fmt.Println("Application ID url param and 'userId' query param must be provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID url param and 'userId' query param must be provided"})
		return
	}

	foundApplication, err := applicationController.findApplicationByIDUseCase.Execute(commands.FindApplicationByID{
		ApplicationID: applicationID,
		QueryByUserID: userID,
	})
	if err != nil {
		fmt.Println("Error while finding application by ID: ", err)
		if _, ok := err.(*errors.UnauthorizedToAccessNamespaceError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deleteApplication := commands.DeleteApplication{
		ID:     applicationID,
		UserID: userID,
	}
	deletedApplication, err := applicationController.deleteApplicationUseCase.Execute(deleteApplication)
	if err != nil {
		fmt.Println("Error while deleting application: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unapplyApplication := commands.UnapplyApplication{
		Name:      foundApplication.Name,
		Namespace: foundApplication.Namespace.Name,
	}
	err = applicationController.undeployApplicationUseCase.Execute(unapplyApplication)
	if err != nil {
		fmt.Println("Error while undeploying application: ", err)
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

	applicationID := c.Param("id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID url param must be provided"})
		return
	}
	queryBy := c.Query("userId")
	if queryBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query param must be provided"})
		return
	}

	application, err := applicationController.findApplicationByIDUseCase.Execute(commands.FindApplicationByID{
		ApplicationID: applicationID,
		QueryByUserID: queryBy,
	})
	if err != nil {
		if _, ok := err.(*errors.UnauthorizedToAccessNamespaceError); ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	applicationNamespace := application.Namespace.Name
	applicationName := application.Name

	getApplicationMetrics := commands.GetApplicationMetrics{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	metrics, err := applicationController.getApplicationMetricsUseCase.Execute(getApplicationMetrics)
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

	applicationID := c.Param("id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID url param must be provided"})
		return
	}
	queryBy := c.Query("userId")
	if queryBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query param must be provided"})
		return
	}

	application, err := applicationController.findApplicationByIDUseCase.Execute(commands.FindApplicationByID{
		ApplicationID: applicationID,
		QueryByUserID: queryBy,
	})
	if err != nil {
		if _, ok := err.(*errors.UnauthorizedToAccessNamespaceError); ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	applicationNamespace := application.Namespace.Name
	applicationName := application.Name
	if applicationNamespace == "" || applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})

		return
	}

	getApplicationLogs := commands.GetApplicationLogs{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	logs, err := applicationController.getApplicationLogsUseCase.Execute(getApplicationLogs)
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

	applicationID := c.Param("id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID url param must be provided"})
		return
	}
	queryBy := c.Query("userId")
	if queryBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query param must be provided"})
		return
	}

	application, err := applicationController.findApplicationByIDUseCase.Execute(commands.FindApplicationByID{
		ApplicationID: applicationID,
		QueryByUserID: queryBy,
	})
	if err != nil {
		if _, ok := err.(*errors.UnauthorizedToAccessNamespaceError); ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	applicationNamespace := application.Namespace.Name
	applicationName := application.Name
	if applicationNamespace == "" || applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace and name must be provided"})
		return
	}

	getApplicationStatus := commands.GetApplicationStatus{
		Name:      applicationName,
		Namespace: applicationNamespace,
	}
	status, err := applicationController.getApplicationStatusUseCase.Execute(getApplicationStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
