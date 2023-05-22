package namespaces

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/repositories"
	"fmt"
)

type FindNamespaceByIDUseCase struct {
	NamespaceRepository repositories.NamespaceRepository
}

func (findNamespaceByIDUseCase FindNamespaceByIDUseCase) Execute(id string) (*domain.Namespace, error) {
	namespace, err := findNamespaceByIDUseCase.NamespaceRepository.FindByID(id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return namespace, nil
}
