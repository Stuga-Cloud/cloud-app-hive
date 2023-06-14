package applications

import (
	"cloud-app-hive/controllers/errors"
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type DeleteApplicationUseCase struct {
	ApplicationRepository repositories.ApplicationRepository
}

func (deleteApplicationUseCase DeleteApplicationUseCase) Execute(deleteApplication commands.DeleteApplication) (*domain.Application, error) {
	application, err := deleteApplicationUseCase.ApplicationRepository.FindByID(deleteApplication.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting application by ID: %w", err)
	}

	isAdmin := false
	for _, membership := range application.Namespace.Memberships {
		if membership.UserID == deleteApplication.UserID && membership.Role == domain.RoleAdmin {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return nil, errors.NewUnauthorizedToAccessNamespaceError(application.Namespace.ID, application.Namespace.Name, deleteApplication.UserID)
	}

	deletedApplication, err := deleteApplicationUseCase.ApplicationRepository.Delete(deleteApplication.ID)
	if err != nil {
		return nil, fmt.Errorf("error when deleting application: %w", err)
	}

	return deletedApplication, nil
}
