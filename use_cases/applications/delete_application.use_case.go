package applications

import (
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

	if application.UserID != deleteApplication.UserID {
		return nil, fmt.Errorf("user is not the owner of the application")
	}

	_, err = deleteApplicationUseCase.ApplicationRepository.Delete(deleteApplication.ID)
	if err != nil {
		return nil, fmt.Errorf("error when deleting application: %w", err)
	}

	return application, nil
}
