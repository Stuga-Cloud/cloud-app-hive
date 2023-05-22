package namespaces

import (
	"github.com/gin-gonic/gin"
)

func InitNamespacesRoutes(router *gin.Engine) {
	router.POST("/namespaces", CreateNamespaceController)
	router.GET("/namespaces", FindNamespacesController)
	router.GET("/namespaces/:id", FindNamespaceByIDController)
}
