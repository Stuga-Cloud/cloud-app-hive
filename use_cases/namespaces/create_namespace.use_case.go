package namespaces

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type CreateNamespaceUseCase struct {
	NamespaceRepository repositories.NamespaceRepository
}

func (createNamespaceUseCase CreateNamespaceUseCase) Execute(createNamespace commands.CreateNamespace) (*domain.Namespace, error) {
	foundNamespaceByID, err := createNamespaceUseCase.NamespaceRepository.ExistsByName(createNamespace.Name)
	if err != nil {
		return nil, err
	}
	if foundNamespaceByID == true {
		return nil, fmt.Errorf("namespace %s already exists", createNamespace.Name)
	}

	createdNamespace, err := createNamespaceUseCase.NamespaceRepository.Create(createNamespace)
	if err != nil {
		fmt.Println(fmt.Errorf("error creating namespace (%s): %w", createNamespace.Name, err))
		return nil, err
	}
	if createdNamespace == nil {
		return nil, fmt.Errorf("namespace %s could not be created", createNamespace.Name)
	}

	return createdNamespace, nil
}
