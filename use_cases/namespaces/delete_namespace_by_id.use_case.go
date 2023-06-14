package namespaces

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/repositories"
)

type DeleteNamespaceByIDUseCase struct {
	NamespaceRepository repositories.NamespaceRepository
}

func (deleteNamespaceByIDUseCase DeleteNamespaceByIDUseCase) Execute(id string, userId string) (*domain.Namespace, error) {
	namespace, err := deleteNamespaceByIDUseCase.NamespaceRepository.Delete(id, userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return namespace, nil
}
