package lambdaControllers

import (
	"cloud-app-hive/controllers/applications/dto"
	"cloud-app-hive/use_cases/applications"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateContainerController(c *gin.Context) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println("Recovered in CreateContainerController", r)
	//		c.JSON(500, gin.H{
	//			"error": "Internal Server Error",
	//		})
	//	}
	//}()

	// Authorization
	//authHeader := c.Request.Header.Get("Authorization")
	//
	//if authHeader != "Bearer " + os.Getenv("API_KEY") {
	//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	//	return
	//}

	// Retrieve body CreateApplicationDto
	var createApplicationDto dto.CreateApplicationDto
	if err := c.ShouldBindJSON(&createApplicationDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := dto.ValidateCreateApplicationDto(createApplicationDto)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	appImage := createApplicationDto.Image
	appName := createApplicationDto.Name
	appNamespace := createApplicationDto.Namespace

	// TODO -> Get subdomain from OVH API
	createSubdomainUseCase := applications.CreateSubdomainUseCase{}
	created := createSubdomainUseCase.Execute(appName, appNamespace)
	fmt.Println("Created subdomain: ", created)

	// TODO -> Deploy application
	deployApplicationUseCase := applications.DeployApplicationUseCase{}
	output, err := deployApplicationUseCase.Execute(appImage, appName, appNamespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while deploying": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("App %s deployed", appName),
		"output":   output,
		"cpuUsage": 0,
		"ramUsage": 0,
		"time":     0,
	})
}
