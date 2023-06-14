package namespaces

import (
	"cloud-app-hive/controllers/errors"
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type UpdateNamespaceByIDUseCase struct {
	NamespaceRepository repositories.NamespaceRepository
}

func (UpdateNamespaceUseCase UpdateNamespaceByIDUseCase) Execute(updateNamespace commands.UpdateNamespace, userID string) (*domain.Namespace, error) {
	foundNamespaceByID, err := UpdateNamespaceUseCase.NamespaceRepository.FindByID(updateNamespace.ID)
	if err != nil {
		return nil, err
	}
	if foundNamespaceByID == nil {
		return nil, errors.NewNamespaceNotFoundByIDError(updateNamespace.ID)
	}

	isAdmin := false
	for _, member := range foundNamespaceByID.Memberships {
		if member.UserID == userID && member.Role == domain.RoleAdmin {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return nil, errors.NewUnauthorizedToAccessNamespaceError(foundNamespaceByID.ID, foundNamespaceByID.Name, userID)
	}

	updatedNamespace, err := UpdateNamespaceUseCase.NamespaceRepository.Update(updateNamespace)
	if err != nil {
		fmt.Println(fmt.Errorf("error updating namespace (%s): %w", updateNamespace.ID, err))
		return nil, err
	}
	if updatedNamespace == nil {
		return nil, fmt.Errorf("namespace %s could not be updated", updateNamespace.ID)
	}

	return updatedNamespace, nil
}
