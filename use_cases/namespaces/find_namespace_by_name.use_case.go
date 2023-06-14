package namespaces

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/repositories"
)

type FindNamespaceByNameUseCase struct {
	NamespaceRepository repositories.NamespaceRepository
}

func (findNamespaceByNameUseCase FindNamespaceByNameUseCase) Execute(name string, userId string) (*domain.Namespace, error) {
	namespace, err := findNamespaceByNameUseCase.NamespaceRepository.FindByName(name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	userCanAccessNamespace := false
	for _, membership := range namespace.Memberships {
		if membership.UserID == userId {
			userCanAccessNamespace = true
		}
	}
	if userCanAccessNamespace == false {
		return nil, fmt.Errorf("user %s does not have access to namespace %s", userId, name)
	}
	return namespace, nil
}
