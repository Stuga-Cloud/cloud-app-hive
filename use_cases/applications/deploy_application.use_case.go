package applications

import (
	"cloud-app-hive/domain"
)

type DeployApplicationUseCase struct {
	// All the repositories that the use case needs
	ContainerManagerRepository domain.ContainerManagerRepository
}

func (deployApplicationUseCase DeployApplicationUseCase) Execute(deployApplication domain.DeployApplication) error {
	err := deployApplicationUseCase.ContainerManagerRepository.DeployApplication(deployApplication)
	if err != nil {
		return err
	}
	return nil
}
