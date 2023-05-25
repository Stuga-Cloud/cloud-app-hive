package namespaces

import (
	"net/http"

	"cloud-app-hive/controllers/namespaces/requests"
	"cloud-app-hive/controllers/validators"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/use_cases/namespaces"
	"github.com/gin-gonic/gin"
)

type NamespaceController struct {
	createNamespaceUseCase   namespaces.CreateNamespaceUseCase
	findNamespacesUseCase    namespaces.FindNamespacesUseCase
	findNamespaceByIDUseCase namespaces.FindNamespaceByIDUseCase
}

func NewNamespaceController(
	createNamespaceUseCase namespaces.CreateNamespaceUseCase,
	findNamespacesUseCase namespaces.FindNamespacesUseCase,
	findNamespaceByIDUseCase namespaces.FindNamespaceByIDUseCase,
) NamespaceController {
	return NamespaceController{
		createNamespaceUseCase:   createNamespaceUseCase,
		findNamespacesUseCase:    findNamespacesUseCase,
		findNamespaceByIDUseCase: findNamespaceByIDUseCase,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	foundNamespace, err := namespaceController.findNamespaceByIDUseCase.Execute(namespaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace": foundNamespace,
	})
}
