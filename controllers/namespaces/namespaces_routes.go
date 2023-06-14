package namespaces

import (
	"cloud-app-hive/use_cases/applications"
	"cloud-app-hive/use_cases/namespaces"
	"github.com/gin-gonic/gin"
)

func InitNamespacesRoutes(
	router *gin.RouterGroup,
	createNamespaceUseCase namespaces.CreateNamespaceUseCase,
	findNamespaceByIDUseCase namespaces.FindNamespaceByIDUseCase,
	findNamespacesUseCase namespaces.FindNamespacesUseCase,
	createNamespaceMembershipUseCase namespaces.CreateNamespaceMembershipUseCase,
	removeNamespaceMembershipUseCase namespaces.RemoveNamespaceMembershipUseCase,
	deleteNamespaceByIDUseCase namespaces.DeleteNamespaceByIDUseCase,
	updateNamespaceByIDUseCase namespaces.UpdateNamespaceByIDUseCase,
	fillApplicationsStatusUseCase applications.FillApplicationStatusUseCase,
) {
	namespaceController := NewNamespaceController(
		createNamespaceUseCase,
		findNamespacesUseCase,
		findNamespaceByIDUseCase,
		createNamespaceMembershipUseCase,
		removeNamespaceMembershipUseCase,
		deleteNamespaceByIDUseCase,
		updateNamespaceByIDUseCase,
		fillApplicationsStatusUseCase,
	)

	router.POST("/namespaces", namespaceController.CreateNamespaceController)
	router.GET("/namespaces", namespaceController.FindNamespacesController)
	router.GET("/namespaces/:id", namespaceController.FindNamespaceByIDController)
	router.DELETE("/namespaces/:id", namespaceController.DeleteNamespaceByIDController)
	router.PUT("/namespaces/:id", namespaceController.UpdateNamespaceByIDController)

	router.POST("/namespaces/:id/memberships", namespaceController.AddMemberToNamespaceController)
	router.DELETE("/namespaces/:id/memberships/:userId", namespaceController.RemoveMemberFromNamespaceController)
}
