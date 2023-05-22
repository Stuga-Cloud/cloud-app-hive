package applications

import (
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
	"fmt"
)

type UndeployApplicationUseCase struct {
	ContainerManagerRepository repositories.ContainerManagerRepository
}

func (undeployApplicationUseCase UndeployApplicationUseCase) Execute(applyApplication commands.UnapplyApplication) error {
	err := undeployApplicationUseCase.ContainerManagerRepository.UnapplyApplication(applyApplication)
	if err != nil {
		return fmt.Errorf("error while applying application: %w", err)
	}
	return nil
}
