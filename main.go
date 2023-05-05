package main

import (
	"cloud-app-hive/config"
	"cloud-app-hive/controllers"
	"context"
	"os"

	"github.com/gin-gonic/gin"
)

var _ = context.Background()

func main() {
	config.Init()

	r := gin.Default()
	r = controllers.InitRoutes(r)
	err := r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
