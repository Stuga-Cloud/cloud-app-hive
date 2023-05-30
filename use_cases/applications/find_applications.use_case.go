package applications

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type FindApplicationsUseCase struct {
	ApplicationRepository repositories.ApplicationRepository
}

func (findApplicationsUseCase FindApplicationsUseCase) Execute(findApplications commands.FindApplications) ([]domain.Application, error) {
	applications, err := findApplicationsUseCase.ApplicationRepository.FindApplications(findApplications)
	if err != nil {
		return nil, fmt.Errorf("error while getting applications: %v", err)
	}

	return applications, nil
}
