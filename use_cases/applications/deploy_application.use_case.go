package applications

import (
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
	"fmt"
)

type DeployApplicationUseCase struct {
	ContainerManagerRepository repositories.ContainerManagerRepository
}

func (deployApplicationUseCase DeployApplicationUseCase) Execute(applyApplication commands.ApplyApplication) error {
	err := deployApplicationUseCase.ContainerManagerRepository.ApplyApplication(applyApplication)
	if err != nil {
		return fmt.Errorf("error while applying application: %w", err)
	}
	return nil
}
