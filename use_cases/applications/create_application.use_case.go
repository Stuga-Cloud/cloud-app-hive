package applications

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type CreateApplicationUseCase struct {
	NamespaceRepository   repositories.NamespaceRepository
	ApplicationRepository repositories.ApplicationRepository
}

func (createApplicationUseCase CreateApplicationUseCase) Execute(createApplication commands.CreateApplication) (*domain.Application, *domain.Namespace, error) {
	foundNamespaceByID, err := createApplicationUseCase.NamespaceRepository.FindByID(createApplication.NamespaceID)
	if err != nil {
		return nil, nil, fmt.Errorf("error while finding namespace by id: %w", err)
	}
	if foundNamespaceByID == nil {
		return nil, nil, fmt.Errorf("no namespace found for namespace id %s", createApplication.NamespaceID)
	}
	if foundNamespaceByID.UserID != createApplication.UserID {
		return nil, nil, fmt.Errorf("user %s is not allowed to access namespace %s", createApplication.UserID, createApplication.NamespaceID)
	}

	foundUserAndNamespaceApplications, err := createApplicationUseCase.ApplicationRepository.FindByNamespaceIDAndUserID(createApplication.NamespaceID, createApplication.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("error while finding applications by namespace id and user id: %w", err)
	}
	if foundUserAndNamespaceApplications != nil {
		for _, foundApplication := range foundUserAndNamespaceApplications {
			if foundApplication.Name == createApplication.Name {
				return nil, nil, fmt.Errorf("application %s already exists in namespace %s", createApplication.Name, createApplication.NamespaceID)
			}
		}
	}

	// Create the application
	createdApplication, err := createApplicationUseCase.ApplicationRepository.Create(createApplication)
	if err != nil {
		return nil, nil, err
	}

	return createdApplication, foundNamespaceByID, nil
}
