package namespaces

import (
	"cloud-app-hive/controllers/namespaces/dto"
	"cloud-app-hive/controllers/validators"
	"cloud-app-hive/database"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/repositories"
	"cloud-app-hive/use_cases/namespaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ApiError from https://github.com/go-playground/validator/issues/559
type ApiError struct {
	Param   string
	Message string
}

func CreateNamespaceController(c *gin.Context) {
	if validators.ValidateAuthorizationToken(c) == false {
		validators.Unauthorized(c)
		return
	}

	if validators.ValidateBodyIsNotNullNorEmpty(c) == false {
		validators.BodyIsNullOrEmptyResponse(c)
		return
	}

	var createNamespaceDto dto.CreateNamespaceDto
	if err := c.ShouldBindJSON(&createNamespaceDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	err := dto.ValidateCreateNamespaceDto(createNamespaceDto)
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
	createNamespaceUseCase := namespaces.CreateNamespaceUseCase{
		NamespaceRepository: repositories.GORMNamespaceRepository{
			Database: db,
		},
	}
	createNamespace := commands.CreateNamespace{
		Name:        createNamespaceDto.Name,
		Description: createNamespaceDto.Description,
		UserID:      createNamespaceDto.UserID,
	}
	namespace, err := createNamespaceUseCase.Execute(createNamespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace": namespace,
	})
}

func FindNamespacesController(c *gin.Context) {
	if validators.ValidateAuthorizationToken(c) == false {
		validators.Unauthorized(c)
		return
	}

	// Validate query params
	var findNamespacesDto dto.FindNamespacesDto
	if err := c.ShouldBindQuery(&findNamespacesDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation-errors": err.Error()})
		return
	}

	err := dto.ValidateFindNamespacesDto(findNamespacesDto)
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
	getNamespacesUseCase := namespaces.FindNamespacesUseCase{
		NamespaceRepository: repositories.GORMNamespaceRepository{
			Database: db,
		},
	}
	findNamespaces := commands.FindNamespaces{
		Name:    findNamespacesDto.Name,
		UserID:  findNamespacesDto.UserID,
		Page:    findNamespacesDto.Page,
		PerPage: findNamespacesDto.PerPage,
	}
	foundNamespaces, err := getNamespacesUseCase.Execute(findNamespaces)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"namespaces": foundNamespaces,
	})
}

func FindNamespaceByIDController(c *gin.Context) {
	if validators.ValidateAuthorizationToken(c) == false {
		validators.Unauthorized(c)
		return
	}

	namespaceID := c.Param("id")

	db, err := database.ConnectToDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	findNamespaceByIDUseCase := namespaces.FindNamespaceByIDUseCase{
		NamespaceRepository: repositories.GORMNamespaceRepository{
			Database: db,
		},
	}
	foundNamespace, err := findNamespaceByIDUseCase.Execute(namespaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace": foundNamespace,
	})
}
