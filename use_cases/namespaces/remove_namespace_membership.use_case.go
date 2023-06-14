package namespaces

import (
	"cloud-app-hive/controllers/errors"
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type RemoveNamespaceMembershipUseCase struct {
	NamespaceRepository           repositories.NamespaceRepository
	NamespaceMembershipRepository repositories.NamespaceMembershipRepository
}

func (removeNamespaceMembershipUseCase RemoveNamespaceMembershipUseCase) Execute(removeNamespaceMembership commands.RemoveNamespaceMembership) (*domain.NamespaceMembership, error) {
	foundNamespaceByID, err := removeNamespaceMembershipUseCase.NamespaceRepository.FindByID(removeNamespaceMembership.NamespaceID)
	if err != nil {
		return nil, err
	}

	foundNamespaceMembershipByID, err := removeNamespaceMembershipUseCase.NamespaceMembershipRepository.ExistsByNamespaceIDAndUserID(removeNamespaceMembership.NamespaceID, removeNamespaceMembership.UserID)
	if err != nil {
		return nil, err
	}
	if foundNamespaceMembershipByID == false {
		return nil, fmt.Errorf("namespace membership of user %s in namespace %s does not exist", removeNamespaceMembership.UserID, removeNamespaceMembership.NamespaceID)
	}

	isUserThatRemovesNamespaceMembershipAdminInNamespace, err := removeNamespaceMembershipUseCase.NamespaceMembershipRepository.IsAdminInNamespace(removeNamespaceMembership.RemovedBy, removeNamespaceMembership.NamespaceID)
	if err != nil {
		return nil, err
	}
	if isUserThatRemovesNamespaceMembershipAdminInNamespace == false {
		return nil, errors.NewNotAdminInNamespaceError(removeNamespaceMembership.NamespaceID, removeNamespaceMembership.RemovedBy)
	}

	isAdminInNamespace, err := removeNamespaceMembershipUseCase.NamespaceMembershipRepository.IsAdminInNamespace(removeNamespaceMembership.UserID, removeNamespaceMembership.NamespaceID)
	if err != nil {
		return nil, err
	}
	// If the user that is being removed is an admin in the namespace, the user that removes the user must be the creator of the namespace
	if isAdminInNamespace == true && foundNamespaceByID.UserID != removeNamespaceMembership.RemovedBy {
		return nil, errors.NewUnauthorizedToRemoveAdminFromNamespaceError(removeNamespaceMembership.NamespaceID, removeNamespaceMembership.UserID, removeNamespaceMembership.RemovedBy)
	}

	removeByNamespaceIDAndUserID, err := removeNamespaceMembershipUseCase.NamespaceMembershipRepository.RemoveByNamespaceIDAndUserID(removeNamespaceMembership.NamespaceID, removeNamespaceMembership.UserID)
	if err != nil {
		fmt.Println(fmt.Errorf("error removing namespace membership (%s): %w", removeNamespaceMembership.NamespaceID, err))
		return nil, err
	}
	if removeByNamespaceIDAndUserID == nil {
		return nil, fmt.Errorf("namespace membership %s could not be removed", removeNamespaceMembership.NamespaceID)
	}

	return removeByNamespaceIDAndUserID, nil
}
