package applications

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
	"fmt"
)

type UpdateApplicationUseCase struct {
	ApplicationRepository repositories.ApplicationRepository
	NamespaceRepository   repositories.NamespaceRepository
}

func (createApplicationUseCase UpdateApplicationUseCase) Execute(applicationID string, updateApplication commands.UpdateApplication) (*domain.Application, *domain.Namespace, error) {
	foundApplicationByID, err := createApplicationUseCase.ApplicationRepository.FindByID(applicationID)
	if err != nil {
		return nil, nil, fmt.Errorf("error while finding application by id: %w", err)
	}
	if foundApplicationByID == nil {
		return nil, nil, fmt.Errorf("no application found for application id %s", applicationID)
	}
	if foundApplicationByID.UserID != updateApplication.UserID {
		return nil, nil, fmt.Errorf("user %s is not allowed to access application %s", updateApplication.UserID, applicationID)
	}

	updatedApplication, err := createApplicationUseCase.ApplicationRepository.Update(applicationID, updateApplication)
	if err != nil {
		return nil, nil, err
	}

	return updatedApplication, &updatedApplication.Namespace, nil
}
