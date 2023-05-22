package main

import (
	"cloud-app-hive/config"
	"cloud-app-hive/controllers"
	"cloud-app-hive/database"
	"context"
	"os"

	"github.com/gin-gonic/gin"
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
	r = controllers.InitRoutes(r)

	err = r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
