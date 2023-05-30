package applications

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/repositories"
)

type FindManualScalingApplicationsUseCase struct {
	ApplicationRepository repositories.ApplicationRepository
}

func (findManualScalingApplicationsUseCase FindManualScalingApplicationsUseCase) Execute() ([]domain.Application, error) {
	applications, err := findManualScalingApplicationsUseCase.ApplicationRepository.FindManualScalingApplications()
	if err != nil {
		return nil, fmt.Errorf("error while getting manual scaling applications: %v", err)
	}
	return applications, nil
}
