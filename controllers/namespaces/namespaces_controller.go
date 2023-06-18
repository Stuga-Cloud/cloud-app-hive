package namespaces

import (
	"cloud-app-hive/controllers/errors"
	"cloud-app-hive/use_cases/applications"
	"fmt"
	"net/http"

	"cloud-app-hive/controllers/namespaces/requests"
	"cloud-app-hive/controllers/validators"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/use_cases/namespaces"
	"github.com/gin-gonic/gin"
)

type NamespaceController struct {
	createNamespaceUseCase           namespaces.CreateNamespaceUseCase
	findNamespacesUseCase            namespaces.FindNamespacesUseCase
	findNamespaceByIDUseCase         namespaces.FindNamespaceByIDUseCase
	findNamespaceByName              namespaces.FindNamespaceByNameUseCase
	createNamespaceMembershipUseCase namespaces.CreateNamespaceMembershipUseCase
	removeNamespaceMembershipUseCase namespaces.RemoveNamespaceMembershipUseCase
	deleteNamespaceByIDUseCase       namespaces.DeleteNamespaceByIDUseCase
	updateNamespaceByIDUseCase       namespaces.UpdateNamespaceByIDUseCase
	fillApplicationsStatusUseCase    applications.FillApplicationStatusUseCase
}

func NewNamespaceController(
	createNamespaceUseCase namespaces.CreateNamespaceUseCase,
	findNamespacesUseCase namespaces.FindNamespacesUseCase,
	findNamespaceByIDUseCase namespaces.FindNamespaceByIDUseCase,
	createNamespaceMembershipUseCase namespaces.CreateNamespaceMembershipUseCase,
	removeNamespaceMembershipUseCase namespaces.RemoveNamespaceMembershipUseCase,
	deleteNamespaceByIDUseCase namespaces.DeleteNamespaceByIDUseCase,
	updateNamespaceByIDUseCase namespaces.UpdateNamespaceByIDUseCase,
	fillApplicationsStatusUseCase applications.FillApplicationStatusUseCase,
) NamespaceController {
	return NamespaceController{
		createNamespaceUseCase:           createNamespaceUseCase,
		findNamespacesUseCase:            findNamespacesUseCase,
		findNamespaceByIDUseCase:         findNamespaceByIDUseCase,
		createNamespaceMembershipUseCase: createNamespaceMembershipUseCase,
		removeNamespaceMembershipUseCase: removeNamespaceMembershipUseCase,
		deleteNamespaceByIDUseCase:       deleteNamespaceByIDUseCase,
		updateNamespaceByIDUseCase:       updateNamespaceByIDUseCase,
		fillApplicationsStatusUseCase:    fillApplicationsStatusUseCase,
	}
}

func (namespaceController NamespaceController) CreateNamespaceController(c *gin.Context) {
	if !validators.ValidateAuthorizationToken(c) {
		validators.Unauthorized(c)
		return
	}

	if !validators.ValidateBodyIsNotNullNorEmpty(c) {
		validators.BodyIsNullOrEmptyResponse(c)
		return
	}

	var createNamespaceRequest requests.CreateNamespaceRequest
	if err := c.ShouldBindJSON(&createNamespaceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	err := requests.ValidateCreateNamespaceRequest(createNamespaceRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"validation-errors": err.Error(),
		})
		return
	}

	createNamespace := commands.CreateNamespace{
		Name:        createNamespaceRequest.Name,
		Description: createNamespaceRequest.Description,
		UserID:      createNamespaceRequest.UserID,
	}
	namespace, err := namespaceController.createNamespaceUseCase.Execute(createNamespace)
	if err != nil {
		if _, ok := err.(*errors.NamespaceWithNameAlreadyExistError); ok {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace": namespace,
	})
}

func (namespaceController NamespaceController) FindNamespacesController(c *gin.Context) {
	if !validators.ValidateAuthorizationToken(c) {
		validators.Unauthorized(c)
		return
	}

	// Validate query params
	var findNamespacesRequest requests.FindNamespacesRequest
	if err := c.ShouldBindQuery(&findNamespacesRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	err := requests.ValidateFindNamespacesRequest(findNamespacesRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"validation-errors": err.Error(),
		})
		return
	}

	findNamespaces := commands.FindNamespaces{
		Name:    findNamespacesRequest.Name,
		UserID:  findNamespacesRequest.UserID,
		Page:    findNamespacesRequest.Page,
		PerPage: findNamespacesRequest.PerPage,
	}
	foundNamespaces, err := namespaceController.findNamespacesUseCase.Execute(findNamespaces)
	if err != nil {
		if _, ok := err.(*errors.UnauthorizedToAccessNamespaceError); ok {
			c.JSON(http.StatusForbidden, errors.NewApiError(
				http.StatusForbidden,
				"user_unauthorized_to_access_namespace",
				fmt.Sprintf("The user with ID %s is not authorized to access the namespace with name %s", findNamespaces.UserID, *findNamespaces.Name),
				"If user should be authorized to access the namespace, please contact the namespace admin(s) to grant access to the user",
				c,
				map[string]interface{}{
					"query_params": c.Request.URL.Query(),
				},
			))
			return
		}
		if _, ok := err.(*errors.NamespaceNotFoundByNameError); ok {
			c.JSON(http.StatusNotFound, errors.NewApiError(
				http.StatusNotFound,
				"namespace_not_found",
				fmt.Sprintf("The namespace with name %s was not found", *findNamespaces.Name),
				"Please try again with a different namespace name",
				c,
				map[string]interface{}{
					"query_params": c.Request.URL.Query(),
				},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, errors.NewApiError(
			http.StatusInternalServerError,
			"internal_server_error",
			"An internal server error occurred while trying to find namespaces",
			"Please try again later",
			c,
			map[string]interface{}{
				"query_params": c.Request.URL.Query(),
			},
		))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"namespaces": foundNamespaces,
	})
}

func (namespaceController NamespaceController) FindNamespaceByIDController(c *gin.Context) {
	if !validators.ValidateAuthorizationToken(c) {
		validators.Unauthorized(c)
		return
	}

	namespaceID := c.Param("id")
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query param is required"})
		return
	}
	foundNamespace, err := namespaceController.findNamespaceByIDUseCase.Execute(namespaceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO Fill applications status by querying it from the container manager
	foundNamespace.Applications, err = namespaceController.fillApplicationsStatusUseCase.Execute(foundNamespace.Name, foundNamespace.Applications)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace": foundNamespace,
	})
}

func (namespaceController NamespaceController) DeleteNamespaceByIDController(c *gin.Context) {
	if !validators.ValidateAuthorizationToken(c) {
		validators.Unauthorized(c)
		return
	}

	namespaceID := c.Param("id")

	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query param is required"})
		return
	}

	namespace, err := namespaceController.deleteNamespaceByIDUseCase.Execute(namespaceID, userID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace": namespace,
	})
}

func (namespaceController NamespaceController) AddMemberToNamespaceController(c *gin.Context) {
	if !validators.ValidateAuthorizationToken(c) {
		validators.Unauthorized(c)
		return
	}

	namespaceID := c.Param("id")
	var addMemberToNamespaceRequest requests.AddMemberToNamespaceRequest
	if err := c.ShouldBindJSON(&addMemberToNamespaceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	createNamespaceMembership := commands.CreateNamespaceMembership{
		NamespaceID: namespaceID,
		UserID:      addMemberToNamespaceRequest.UserID,
		Role:        addMemberToNamespaceRequest.Role,
		AddedBy:     addMemberToNamespaceRequest.AddedBy,
	}
	namespaceMembership, err := namespaceController.createNamespaceMembershipUseCase.Execute(createNamespaceMembership)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace_membership": namespaceMembership,
	})
}

func (namespaceController NamespaceController) RemoveMemberFromNamespaceController(c *gin.Context) {
	if !validators.ValidateAuthorizationToken(c) {
		validators.Unauthorized(c)
		return
	}

	namespaceID := c.Param("id")
	userID := c.Param("userId")
	removedBy := c.Query("removedBy")
	if removedBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "removedBy query param is required"})
		return
	}

	namespaceMembership, err := namespaceController.removeNamespaceMembershipUseCase.Execute(commands.RemoveNamespaceMembership{
		NamespaceID: namespaceID,
		UserID:      userID,
		RemovedBy:   removedBy,
	})
	if err != nil {
		fmt.Printf("Error while removing member from namespace: %s", err.Error())

		status := http.StatusInternalServerError
		errorName := "internal_server_error"
		errorMessage := "An internal server error occurred while trying to remove the member from the namespace"
		description := "Please try again later or contact support"
		contextualInformation := map[string]interface{}{
			"queryParams":  c.Request.URL.Query(),
			"namespace_id": namespaceID,
			"userId":       userID,
			"removedBy":    removedBy,
		}
		if _, ok := err.(*errors.UnauthorizedToAccessNamespaceError); ok {
			status = http.StatusForbidden
			errorName = "unauthorized_to_access_namespace"
			errorMessage = fmt.Sprintf("The user with id %s is not authorized to access the namespace with id %s", userID, namespaceID)
			description = "Please try again with a different namespace id"
		}
		if _, ok := err.(*errors.NotAdminInNamespaceError); ok {
			status = http.StatusForbidden
			errorName = "not_admin_in_namespace"
			errorMessage = fmt.Sprintf("The user with id %s is not an admin in the namespace with id %s", removedBy, namespaceID)
			description = "Maybe try again with a different namespace id or ask admin to remove this user from the namespace"
		}
		if _, ok := err.(*errors.UnauthorizedToRemoveAdminFromNamespaceError); ok {
			status = http.StatusForbidden
			errorName = "unauthorized_to_remove_admin_from_namespace"
			errorMessage = fmt.Sprintf("The user with id %s is not authorized to remove an admin from the namespace with id %s", removedBy, namespaceID)
			description = "Please try again with a different namespace id. Because you cannot remove another admin from the namespace"
		}
		c.JSON(http.StatusInternalServerError, errors.NewApiError(
			status,
			errorName,
			errorMessage,
			description,
			c,
			contextualInformation,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace_membership": namespaceMembership,
	})
}

func (namespaceController NamespaceController) UpdateNamespaceByIDController(c *gin.Context) {
	if !validators.ValidateAuthorizationToken(c) {
		validators.Unauthorized(c)
		return
	}

	namespaceID := c.Param("id")
	userID := c.Query("userId")
	var updateNamespaceRequest requests.UpdateNamespaceByIDRequest
	if err := c.ShouldBindJSON(&updateNamespaceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	updateNamespace := commands.UpdateNamespace{
		ID:          namespaceID,
		UserID:      updateNamespaceRequest.UserID,
		Description: updateNamespaceRequest.Description,
	}
	namespace, err := namespaceController.updateNamespaceByIDUseCase.Execute(updateNamespace, userID)
	if err != nil {
		if _, ok := err.(*errors.NamespaceNotFoundByIDError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace": namespace,
	})
}
