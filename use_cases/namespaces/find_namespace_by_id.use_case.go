package namespaces

import (
	"cloud-app-hive/controllers/errors"
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/repositories"
)

type FindNamespaceByIDUseCase struct {
	NamespaceRepository repositories.NamespaceRepository
}

func (findNamespaceByIDUseCase FindNamespaceByIDUseCase) Execute(id, userId string) (*domain.Namespace, error) {
	namespace, err := findNamespaceByIDUseCase.NamespaceRepository.FindByID(id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if namespace == nil {
		return nil, fmt.Errorf("namespace not found with ID %s", id)
	}

	isAdmin := false
	for _, member := range namespace.Memberships {
		if member.UserID == userId && member.Role == domain.RoleAdmin {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return nil, errors.NewUnauthorizedToAccessNamespaceError(namespace.ID, namespace.Name, userId)
	}

	return namespace, nil
}
