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
	createNamespaceMembershipUseCase namespaces.CreateNamespaceMembershipUseCase,
	deleteNamespaceByIDUseCase namespaces.DeleteNamespaceByIDUseCase,
	updateNamespaceByIDUseCase namespaces.UpdateNamespaceByIDUseCase,
) {
	namespaceController := NewNamespaceController(
		createNamespaceUseCase,
		findNamespacesUseCase,
		findNamespaceByIDUseCase,
		createNamespaceMembershipUseCase,
		deleteNamespaceByIDUseCase,
		updateNamespaceByIDUseCase,
	)

	router.POST("/namespaces", namespaceController.CreateNamespaceController)
	router.GET("/namespaces", namespaceController.FindNamespacesController)
	router.GET("/namespaces/:id", namespaceController.FindNamespaceByIDController)
	router.DELETE("/namespaces/:id", namespaceController.DeleteNamespaceByIDController)
	router.PUT("/namespaces/:id", namespaceController.UpdateNamespaceByIDController)

	router.POST("/namespaces/:id/memberships", namespaceController.AddMemberToNamespaceController)
}
