package applications

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type GetApplicationStatusUseCase struct {
	ContainerManagerRepository repositories.ContainerManagerRepository
}

func (getApplicationStatusUseCase GetApplicationStatusUseCase) Execute(application commands.GetApplicationStatus) (*domain.ApplicationStatus, error) {
	status, err := getApplicationStatusUseCase.ContainerManagerRepository.GetApplicationStatus(application)
	if err != nil {
		return nil, err
	}
	return status, nil
}
