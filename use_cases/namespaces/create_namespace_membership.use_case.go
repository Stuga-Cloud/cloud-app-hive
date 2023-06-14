package namespaces

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type CreateNamespaceMembershipUseCase struct {
	NamespaceMembershipRepository repositories.NamespaceMembershipRepository
}

func (createNamespaceMembershipUseCase CreateNamespaceMembershipUseCase) Execute(createNamespaceMembership commands.CreateNamespaceMembership) (*domain.NamespaceMembership, error) {
	foundNamespaceMembershipByID, err := createNamespaceMembershipUseCase.NamespaceMembershipRepository.ExistsByNamespaceIDAndUserID(createNamespaceMembership.NamespaceID, createNamespaceMembership.UserID)
	if err != nil {
		return nil, err
	}
	if foundNamespaceMembershipByID == true {
		return nil, fmt.Errorf("namespace membership %s already exists", createNamespaceMembership.NamespaceID)
	}

	createdNamespaceMembership, err := createNamespaceMembershipUseCase.NamespaceMembershipRepository.Create(createNamespaceMembership)
	if err != nil {
		fmt.Println(fmt.Errorf("error creating namespace membership (%s): %w", createNamespaceMembership.NamespaceID, err))
		return nil, err
	}
	if createdNamespaceMembership == nil {
		return nil, fmt.Errorf("namespace membership %s could not be created", createNamespaceMembership.NamespaceID)
	}

	return createdNamespaceMembership, nil
}
