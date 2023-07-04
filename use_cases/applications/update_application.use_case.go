package applications

import (
	"fmt"

	"cloud-app-hive/controllers/errors"
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type UpdateApplicationUseCase struct {
	ApplicationRepository repositories.ApplicationRepository
	NamespaceRepository   repositories.NamespaceRepository
}

func (createApplicationUseCase UpdateApplicationUseCase) Execute(applicationID string, updateApplication commands.UpdateApplication, byUserID string) (*domain.Application, *domain.Namespace, error) {
	foundApplicationByID, err := createApplicationUseCase.ApplicationRepository.FindByID(applicationID)
	if err != nil {
		return nil, nil, fmt.Errorf("error while finding application by id: %w", err)
	}
	if foundApplicationByID == nil {
		return nil, nil, fmt.Errorf("no application found for application id %s", applicationID)
	}

	isAdmin := false
	for _, member := range foundApplicationByID.Namespace.Memberships {
		if member.UserID == byUserID && member.Role == domain.RoleAdmin {
			isAdmin = true
			break
		}
	}
	isAppOwner := foundApplicationByID.UserID == byUserID
	if !isAdmin && !isAppOwner {
		return nil, nil, errors.NewUnauthorizedToAccessNamespaceError(foundApplicationByID.Namespace.ID, foundApplicationByID.Namespace.Name, byUserID)
	}

	updatedApplication, err := createApplicationUseCase.ApplicationRepository.Update(applicationID, updateApplication)
	if err != nil {
		return nil, nil, err
	}

	return updatedApplication, &updatedApplication.Namespace, nil
}
