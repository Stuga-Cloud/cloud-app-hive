package namespaces

import (
	"cloud-app-hive/controllers/errors"
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type FindNamespacesUseCase struct {
	NamespaceRepository repositories.NamespaceRepository
}

func (findNamespacesUseCase FindNamespacesUseCase) Execute(findNamespaces commands.FindNamespaces) ([]domain.Namespace, error) {
	namespaces, err := findNamespacesUseCase.NamespaceRepository.Find(findNamespaces)
	if err != nil {
		return nil, fmt.Errorf("error finding namespaces: %w", err)
	}
	if namespaces == nil || len(namespaces) == 0 {
		return nil, errors.NewNamespaceNotFoundByNameError(*findNamespaces.Name)
	}
	// Check that user has access to all namespaces
	for _, namespace := range namespaces {
		for _, namespaceMembership := range namespace.Memberships {
			if namespaceMembership.UserID != findNamespaces.UserID {
				return nil, errors.NewUnauthorizedToAccessNamespaceError(namespace.ID, namespace.Name, findNamespaces.UserID)
			}
		}
	}
	return namespaces, nil
}
