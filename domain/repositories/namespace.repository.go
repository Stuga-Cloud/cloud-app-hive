package repositories

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
)

// NamespaceRepository is an interface that represents a repository of namespaces
type NamespaceRepository interface {
	// FindByID returns a namespace by its ID
	FindByID(id string) (*domain.Namespace, error)
	// ExistsByName
	ExistsByName(name string) (bool, error)
	// FindByName returns a namespace by its name
	FindByName(name string) (*domain.Namespace, error)
	// Find returns a list of namespaces
	Find(findNamespaces commands.FindNamespaces) ([]domain.Namespace, error)
	// Create creates a new namespace
	Create(namespace commands.CreateNamespace) (*domain.Namespace, error)
	// Delete deletes a namespace
	Delete(id string, userId string) (*domain.Namespace, error)
	// Update updates a namespace
	Update(namespace commands.UpdateNamespace) (*domain.Namespace, error)
}
