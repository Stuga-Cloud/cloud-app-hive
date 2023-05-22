package main

import (
	"cloud-app-hive/config"
	"cloud-app-hive/controllers"
	"cloud-app-hive/database"
	"cloud-app-hive/docs"
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

var _ = context.Background()

func main() {
	config.Init()

	db, err := database.ConnectToDatabase()
	if err != nil {
		panic(err)
	}
	err = database.MigrateDatabase(db)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	docs.SwaggerInfo.Title = "Cloud App Hive API"
	docs.SwaggerInfo.Description = "This API is used to manage applications in the cloud"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r = controllers.InitRoutes(r)

	err = r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
