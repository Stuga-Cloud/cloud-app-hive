package repositories

import (
	"cloud-app-hive/controllers/errors"
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type GORMNamespaceMembershipRepository struct {
	Database *gorm.DB
}

// ExistsByNamespaceIDAndUserID returns true if a namespace membership exists by namespace ID and user ID
func (r GORMNamespaceMembershipRepository) ExistsByNamespaceIDAndUserID(namespaceID string, userID string) (bool, error) {
	namespaceMembership := domain.NamespaceMembership{}
	result := r.Database.Find(&namespaceMembership, domain.NamespaceMembership{
		NamespaceID: namespaceID,
		UserID:      userID,
	}).Limit(1)
	if result.Error != nil {
		return false, fmt.Errorf("error checking if namespace membership exists: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}

// FindByUserID returns a list of namespaces by user ID
func (r GORMNamespaceMembershipRepository) FindByUserID(userID string) ([]domain.NamespaceMembership, error) {
	namespaceMemberships := []domain.NamespaceMembership{}
	result := r.Database.Find(&namespaceMemberships, domain.NamespaceMembership{
		UserID: userID,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error finding namespace memberships: %w", result.Error)
	}
	return namespaceMemberships, nil
}

// FindByNamespaceID returns a list of namespaces by namespace ID
func (r GORMNamespaceMembershipRepository) FindByNamespaceID(namespaceID string) ([]domain.NamespaceMembership, error) {
	namespaceMemberships := []domain.NamespaceMembership{}
	result := r.Database.Find(&namespaceMemberships, domain.NamespaceMembership{
		NamespaceID: namespaceID,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error finding namespace memberships: %w", result.Error)
	}
	return namespaceMemberships, nil
}

// FindByUserIDAndNamespaceID returns a namespace membership by user ID and namespace ID
func (r GORMNamespaceMembershipRepository) FindByUserIDAndNamespaceID(userID string, namespaceID string) (*domain.NamespaceMembership, error) {
	namespaceMembership := domain.NamespaceMembership{}
	result := r.Database.Find(&namespaceMembership, domain.NamespaceMembership{
		UserID:      userID,
		NamespaceID: namespaceID,
	}).Limit(1)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding namespace membership: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("namespace membership not found with user ID %s and namespace ID %s", userID, namespaceID)
	}
	return &namespaceMembership, nil
}

// Create creates a new namespace membership
func (r GORMNamespaceMembershipRepository) Create(createNamespaceMembership commands.CreateNamespaceMembership) (*domain.NamespaceMembership, error) {
	namespace := domain.Namespace{}
	result := r.Database.Preload("Memberships").Find(&namespace, domain.Namespace{
		ID: createNamespaceMembership.NamespaceID,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error finding namespace: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("namespace not found with ID %s", createNamespaceMembership.NamespaceID)
	}

	isAuthorized := false
	for _, namespaceMembership := range namespace.Memberships {
		if namespaceMembership.UserID == createNamespaceMembership.AddedBy && namespaceMembership.Role == domain.RoleAdmin {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		return nil, fmt.Errorf("user %s is not authorized to add users to namespace %s", createNamespaceMembership.AddedBy, createNamespaceMembership.NamespaceID)
	}

	namespaceMembershipModel := domain.NamespaceMembership{
		ID:          uuid.New().String(),
		NamespaceID: createNamespaceMembership.NamespaceID,
		UserID:      createNamespaceMembership.UserID,
		Role:        createNamespaceMembership.Role,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	resultCreation := r.Database.Create(&namespaceMembershipModel)
	if resultCreation.Error != nil {
		return nil, fmt.Errorf("error creating namespace membership: %w", resultCreation.Error)
	}
	return &namespaceMembershipModel, nil
}

// Delete deletes a namespace membership
func (r GORMNamespaceMembershipRepository) Delete(namespaceMembershipID string) (*domain.NamespaceMembership, error) {
	namespaceMembership := domain.NamespaceMembership{}
	result := r.Database.Delete(&namespaceMembership, namespaceMembershipID)
	if result.Error != nil {
		return nil, fmt.Errorf("error deleting namespace membership: %w", result.Error)
	}
	return &namespaceMembership, nil
}

// RemoveByNamespaceIDAndUserID removes a namespace membership by namespace ID and user ID
func (r GORMNamespaceMembershipRepository) RemoveByNamespaceIDAndUserID(namespaceID string, userID string) (*domain.NamespaceMembership, error) {
	namespaceMembership := domain.NamespaceMembership{}

	result := r.Database.Delete(&namespaceMembership, domain.NamespaceMembership{
		NamespaceID: namespaceID,
		UserID:      userID,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error deleting namespace membership: %w", result.Error)
	}
	return &namespaceMembership, nil
}

// IsAdminInNamespace returns true if the user is an admin in the namespace
func (r GORMNamespaceMembershipRepository) IsAdminInNamespace(userID string, namespaceID string) (bool, error) {
	namespaceMembership := domain.NamespaceMembership{}
	result := r.Database.Find(&namespaceMembership, domain.NamespaceMembership{
		UserID:      userID,
		NamespaceID: namespaceID,
	}).Limit(1)
	if result.Error != nil {
		return false, fmt.Errorf("error finding namespace membership: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return false, errors.NewUnauthorizedToAccessNamespaceError(namespaceID, "", userID)
	}
	return namespaceMembership.Role == domain.RoleAdmin, nil
}

// Update updates a namespace membership
func (r GORMNamespaceMembershipRepository) Update(namespaceMembership commands.UpdateNamespaceMembership) (*domain.NamespaceMembership, error) {
	namespaceMembershipModel := domain.NamespaceMembership{}
	result := r.Database.Model(&namespaceMembershipModel).Where(domain.NamespaceMembership{
		UserID:      namespaceMembership.UserID,
		NamespaceID: namespaceMembership.NamespaceID,
	}).Updates(domain.NamespaceMembership{
		Role:      namespaceMembership.Role,
		UpdatedAt: time.Now().UTC(),
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error updating namespace membership: %w", result.Error)
	}
	return &namespaceMembershipModel, nil
}
