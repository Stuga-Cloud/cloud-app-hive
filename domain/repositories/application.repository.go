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

	// FindByUserID returns applications by its user ID
	FindByUserID(userID string) ([]domain.Application, error)

	// FindByNamespaceIDAndUserID returns an application by namespace ID
	FindByNamespaceIDAndUserID(namespaceID string) ([]domain.Application, error)

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

	// FindAutoScalingApplications returns all auto scaling applications
	FindAutoScalingApplications() ([]domain.Application, error)

	// HorizontalScaleUp scales up an application horizontally
	HorizontalScaleUp(applicationID string) (*domain.Application, error)

	// HorizontalScaleDown scales down an application horizontally
	HorizontalScaleDown(applicationID string) (*domain.Application, error)

	// // VerticalScaleUp scales up an application vertically
	VerticalScaleUp(applicationID string) (*domain.Application, error)

	// // VerticalScaleDown scales down an application vertically
	// VerticalScaleDown(applicationID string) (*domain.Application, error)
}
