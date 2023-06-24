package applications

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/repositories"
)

type FindAutoScalingApplicationsUseCase struct {
	ApplicationRepository repositories.ApplicationRepository
}

func (findAutoScalingApplicationsUseCase FindAutoScalingApplicationsUseCase) Execute() ([]domain.Application, error) {
	applications, err := findAutoScalingApplicationsUseCase.ApplicationRepository.FindAutoScalingApplications()
	if err != nil {
		return nil, fmt.Errorf("error while getting auto scaling applications: %v", err)
	}
	return applications, nil
}
