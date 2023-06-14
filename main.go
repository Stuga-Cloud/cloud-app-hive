package main

import (
	"cloud-app-hive/repositories"
	"cloud-app-hive/use_cases/applications"
	"cloud-app-hive/use_cases/namespaces"
	"os"

	"cloud-app-hive/config"
	"cloud-app-hive/controllers"
	"cloud-app-hive/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

import (
	"context"

	"cloud-app-hive/docs"
)

var _ = context.Background()

func main() {
	config.InitEnvironmentFile()

	router := gin.Default()
	configApp := cors.DefaultConfig()
	configApp.AllowAllOrigins = true
	configApp.AllowHeaders = []string{"*"}
	configApp.AllowMethods = []string{"*"}
	configApp.MaxAge = 0
	router.Use(cors.New(configApp))

	docs.SwaggerInfo.Title = "Cloud App Hive API"
	docs.SwaggerInfo.Description = "This API is used to manage applications in the cloud"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.GET("/openapi/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	initDependencies(router)

	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		panic(err)
	}
}

func initDependencies(router *gin.Engine) {
	db, err := database.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	if err = database.MigrateDatabase(db); err != nil {
		panic(err)
	}

	// Namespace dependencies
	namespaceRepository := repositories.GORMNamespaceRepository{
		Database: db,
	}
	createNamespaceUseCase := namespaces.CreateNamespaceUseCase{
		NamespaceRepository: namespaceRepository,
	}
	findNamespaceByIDUseCase := namespaces.FindNamespaceByIDUseCase{
		NamespaceRepository: namespaceRepository,
	}
	findNamespacesUseCase := namespaces.FindNamespacesUseCase{
		NamespaceRepository: namespaceRepository,
	}
	deleteNamespaceByIDUseCase := namespaces.DeleteNamespaceByIDUseCase{
		NamespaceRepository: namespaceRepository,
	}
	updateNamespaceByIDUseCase := namespaces.UpdateNamespaceByIDUseCase{
		NamespaceRepository: namespaceRepository,
	}

	// Application dependencies
	applicationRepository := repositories.GORMApplicationRepository{
		Database: db,
	}
	containerManagerRepository := repositories.KubernetesContainerManagerRepository{}
	findApplicationsUseCase := applications.FindApplicationsUseCase{
		ApplicationRepository: applicationRepository,
	}
	createApplicationUseCase := applications.CreateApplicationUseCase{
		ApplicationRepository: applicationRepository,
		NamespaceRepository:   namespaceRepository,
	}
	updateApplicationUseCase := applications.UpdateApplicationUseCase{
		ApplicationRepository: applicationRepository,
		NamespaceRepository:   namespaceRepository,
	}
	deleteApplicationUseCase := applications.DeleteApplicationUseCase{
		ApplicationRepository: applicationRepository,
	}
	deployApplicationUseCase := applications.DeployApplicationUseCase{
		ContainerManagerRepository: containerManagerRepository,
	}
	undeployApplicationUseCase := applications.UndeployApplicationUseCase{
		ContainerManagerRepository: containerManagerRepository,
	}
	getApplicationLogsUseCase := applications.GetApplicationLogsUseCase{
		ContainerManagerRepository: containerManagerRepository,
	}
	getApplicationMetricsUseCase := applications.GetApplicationMetricsUseCase{
		ContainerManagerRepository: containerManagerRepository,
	}
	getApplicationStatusUseCase := applications.GetApplicationStatusUseCase{
		ContainerManagerRepository: containerManagerRepository,
	}
	fillApplicationsStatusUseCase := applications.FillApplicationStatusUseCase{
		ContainerManagerRepository: containerManagerRepository,
	}

	// Namespace membership dependencies
	memoryNamespaceMembershipRepository := repositories.GORMNamespaceMembershipRepository{
		Database: db,
	}
	createNamespaceMembershipUseCase := namespaces.CreateNamespaceMembershipUseCase{
		NamespaceMembershipRepository: memoryNamespaceMembershipRepository,
	}
	removeNamespaceMembershipUseCase := namespaces.RemoveNamespaceMembershipUseCase{
		NamespaceMembershipRepository: memoryNamespaceMembershipRepository,
		NamespaceRepository:           namespaceRepository,
	}

	controllers.InitRoutes(
		router,
		createNamespaceUseCase,
		findNamespaceByIDUseCase,
		findNamespacesUseCase,
		findApplicationsUseCase,
		createApplicationUseCase,
		updateApplicationUseCase,
		deleteApplicationUseCase,
		deployApplicationUseCase,
		undeployApplicationUseCase,
		getApplicationLogsUseCase,
		getApplicationMetricsUseCase,
		getApplicationStatusUseCase,
		fillApplicationsStatusUseCase,
		createNamespaceMembershipUseCase,
		removeNamespaceMembershipUseCase,
		deleteNamespaceByIDUseCase,
		updateNamespaceByIDUseCase,
	)

	//schedulers.InitSchedulers()
}
