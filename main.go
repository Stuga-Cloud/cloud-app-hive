package main

import (
	"os"

	"cloud-app-hive/config"
	"cloud-app-hive/controllers"
	"cloud-app-hive/database"

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
	config.Init()

	db, err := database.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	if err = database.MigrateDatabase(db); err != nil {
		panic(err)
	}

	router := gin.Default()

	docs.SwaggerInfo.Title = "Cloud App Hive API"
	docs.SwaggerInfo.Description = "This API is used to manage applications in the cloud"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.GET("/openapi/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router = controllers.InitRoutes(router)

	err = router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
