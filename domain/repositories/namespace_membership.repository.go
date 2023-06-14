package repositories

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
)

// NamespaceMembershipRepository is an interface that represents a repository of namespace memberships
type NamespaceMembershipRepository interface {
	// ExistsByNamespaceIDAndUserID returns true if a namespace membership exists by namespace ID and user ID
	ExistsByNamespaceIDAndUserID(namespaceID string, userID string) (bool, error)
	// FindByUserID returns a list of namespaces by user ID
	FindByUserID(userID string) ([]domain.NamespaceMembership, error)
	// FindByNamespaceID returns a list of namespaces by namespace ID
	FindByNamespaceID(namespaceID string) ([]domain.NamespaceMembership, error)
	// FindByUserIDAndNamespaceID returns a namespace membership by user ID and namespace ID
	FindByUserIDAndNamespaceID(userID string, namespaceID string) (*domain.NamespaceMembership, error)
	// Create creates a new namespace membership
	Create(namespaceMembership commands.CreateNamespaceMembership) (*domain.NamespaceMembership, error)
	// Delete deletes a namespace membership
	Delete(namespaceMembershipID string) (*domain.NamespaceMembership, error)
	// RemoveByNamespaceIDAndUserID removes a namespace membership by namespace ID and user ID
	RemoveByNamespaceIDAndUserID(userID string, namespaceID string) (*domain.NamespaceMembership, error)
	// IsAdminInNamespace returns true if a user is an admin in a namespace
	IsAdminInNamespace(namespaceID string, userID string) (bool, error)
	// Update updates a namespace membership
	Update(namespaceMembership commands.UpdateNamespaceMembership) (*domain.NamespaceMembership, error)
}
