package applications

import (
	"cloud-app-hive/controllers/errors"
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type FindApplicationByIDUseCase struct {
	ApplicationRepository repositories.ApplicationRepository
}

func (findApplicationByIDUseCase FindApplicationByIDUseCase) Execute(findApplicationByID commands.FindApplicationByID) (*domain.Application, error) {
	application, err := findApplicationByIDUseCase.ApplicationRepository.FindByID(findApplicationByID.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("error while finding application by id: %w", err)
	}
	if application == nil {
		return nil, errors.NewApplicationNotFoundByIDError(findApplicationByID.ApplicationID)
	}

	// Check that user is in namespace memberships
	isAllowed := false
	for _, membership := range application.Namespace.Memberships {
		if membership.UserID == findApplicationByID.QueryByUserID {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return nil, errors.NewUnauthorizedToAccessNamespaceError(
			application.Namespace.ID,
			application.Namespace.Name,
			findApplicationByID.QueryByUserID,
		)
	}

	return application, nil
}
