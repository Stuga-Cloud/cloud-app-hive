package namespaces

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/repositories"
)

type DeleteNamespaceByIDUseCase struct {
	NamespaceRepository        repositories.NamespaceRepository
	ContainerManagerRepository repositories.ContainerManagerRepository
}

func (deleteNamespaceByIDUseCase DeleteNamespaceByIDUseCase) Execute(id string, userId string) (*domain.Namespace, error) {
	foundNamespace, err := deleteNamespaceByIDUseCase.NamespaceRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if foundNamespace == nil {
		return nil, fmt.Errorf("namespace not found with ID %s while deleting", id)
	}

	isAdmin := false
	for _, member := range foundNamespace.Memberships {
		if member.UserID == userId && member.Role == domain.RoleAdmin {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return nil, fmt.Errorf("user %s is not admin of namespace %s, he cannot delete namespace", userId, id)
	}

	err = deleteNamespaceByIDUseCase.ContainerManagerRepository.DeleteNamespace(foundNamespace.Name)
	if err != nil {
		return nil, err
	}

	namespace, err := deleteNamespaceByIDUseCase.NamespaceRepository.Delete(id, userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return namespace, nil
}
