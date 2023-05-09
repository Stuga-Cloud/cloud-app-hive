package applications

import (
	"cloud-app-hive/domain"
)

type DeployApplicationUseCase struct {
	// All the repositories that the use case needs
	ContainerManagerRepository domain.ContainerManagerRepository
}

func (deployApplicationUseCase DeployApplicationUseCase) Execute(deployApplication domain.DeployApplication) (string, error) {
	output, err := deployApplicationUseCase.ContainerManagerRepository.DeployApplication(deployApplication)
	if err != nil {
		return "", err
	}
	println("Deployed FULLY : ", output)

	return "output", nil
}
