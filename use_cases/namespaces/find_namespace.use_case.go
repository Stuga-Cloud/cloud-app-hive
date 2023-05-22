package namespaces

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
	"fmt"
)

type FindNamespacesUseCase struct {
	NamespaceRepository repositories.NamespaceRepository
}

func (findNamespacesUseCase FindNamespacesUseCase) Execute(findNamespaces commands.FindNamespaces) ([]domain.Namespace, error) {
	namespaces, err := findNamespacesUseCase.NamespaceRepository.Find(findNamespaces)
	if err != nil {
		return nil, fmt.Errorf("error finding namespaces: %w", err)
	}
	return namespaces, nil
}
