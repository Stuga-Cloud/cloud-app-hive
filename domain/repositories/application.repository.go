package repositories

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
)

// ApplicationRepository is an interface that represents a repository of applications
type ApplicationRepository interface {
	// FindApplications returns all applications by applying the given filters
	FindApplications(findApplications commands.FindApplications) ([]domain.Application, error)

	// FindByID returns an application by its ID
	FindByID(id string) (*domain.Application, error)

	// FindByUserID returns an application by its user ID
	FindByUserID(userID string) (*domain.Application, error)

	// FindByNamespaceIDAndUserID returns an application by namespace ID and user ID
	FindByNamespaceIDAndUserID(namespaceID string, userID string) ([]domain.Application, error)

	// FindByNamespaceIDAndName returns an application by namespace ID and name
	FindByNamespaceIDAndName(namespaceID string, name string) (*domain.Application, error)

	// Create creates a new application
	Create(application commands.CreateApplication) (*domain.Application, error)

	// Update updates an application
	Update(applicationID string, application commands.UpdateApplication) (*domain.Application, error)

	// Delete deletes an application
	Delete(id string) (*domain.Application, error)

	// FindManualScalingApplications returns all manual scaling applications
	FindManualScalingApplications() ([]domain.Application, error)
}
