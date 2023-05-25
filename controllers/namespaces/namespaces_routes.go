package namespaces

import (
	"cloud-app-hive/use_cases/namespaces"
	"github.com/gin-gonic/gin"
)

func InitNamespacesRoutes(
	router *gin.RouterGroup,
	createNamespaceUseCase namespaces.CreateNamespaceUseCase,
	findNamespaceByIDUseCase namespaces.FindNamespaceByIDUseCase,
	findNamespacesUseCase namespaces.FindNamespacesUseCase,
) {
	namespaceController := NewNamespaceController(
		createNamespaceUseCase,
		findNamespacesUseCase,
		findNamespaceByIDUseCase,
	)

	router.POST("/namespaces", namespaceController.CreateNamespaceController)
	router.GET("/namespaces", namespaceController.FindNamespacesController)
	router.GET("/namespaces/:id", namespaceController.FindNamespaceByIDController)
}
