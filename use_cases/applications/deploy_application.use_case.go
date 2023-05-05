package applications

import (
	"cloud-app-hive/utils"
	"fmt"
)

type DeployApplicationUseCase struct {
	// All the repositories that the use case needs
}

func (deployApplicationUseCase DeployApplicationUseCase) Execute(appImage, appName, appNamespace string) (string, error) {
	clientset, err := utils.ConnectToKubernetesAPI()
	if err != nil {
		panic(err.Error())
	}
	services, err := utils.GetKubernetesServices(clientset, "default")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Services: ", services)

	return "output", nil
}
