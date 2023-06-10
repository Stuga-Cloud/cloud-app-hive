package repositories

import (
	"fmt"
	"time"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GORMNamespaceRepository struct {
	Database *gorm.DB
}

// FindByID returns a namespace by its ID
func (r GORMNamespaceRepository) FindByID(id string) (*domain.Namespace, error) {
	app := domain.Namespace{}
	result := r.Database.Find(&app, domain.Namespace{
		ID: id,
	}).Limit(1)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding namespace: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("namespace not found with ID %s", id)
	}
	return &app, nil
}

// ExistsByName returns true if a namespace with the given name exists
func (r GORMNamespaceRepository) ExistsByName(name string) (bool, error) {
	var count int64
	result := r.Database.Model(&domain.Namespace{}).Where(domain.Namespace{
		Name: name,
	}).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("error checking if namespace exists: %w", result.Error)
	}
	return count > 0, nil
}

// FindByName returns a namespace by its name
func (r GORMNamespaceRepository) FindByName(name string) (*domain.Namespace, error) {
	app := domain.Namespace{}
	result := r.Database.Find(&app, domain.Namespace{
		Name: name,
	}).Limit(1)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding namespace: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("namespace not found with name %s", name)
	}
	return &app, nil
}

// Find returns a list of namespaces
func (r GORMNamespaceRepository) Find(findNamespaces commands.FindNamespaces) ([]domain.Namespace, error) {
	var foundNamespaces []domain.Namespace
	result := r.Database.Find(&foundNamespaces, domain.Namespace{
		Name:   findNamespaces.Name,
		UserID: findNamespaces.UserID,
	}).Limit(findNamespaces.PerPage).Offset(findNamespaces.PerPage * (findNamespaces.Page - 1))
	if result.Error != nil {
		return nil, fmt.Errorf("error finding namespaces: %w", result.Error)
	}
	return foundNamespaces, nil
}

// Create creates a new namespace
func (r GORMNamespaceRepository) Create(createNamespace commands.CreateNamespace) (*domain.Namespace, error) {
	namespace := domain.Namespace{
		ID:          uuid.New().String(),
		Name:        createNamespace.Name,
		Description: createNamespace.Description,
		UserID:      createNamespace.UserID,
	}
	result := r.Database.Create(&namespace)
	if result.Error != nil {
		return nil, fmt.Errorf("error creating namespace: %w", result.Error)
	}
	return &namespace, nil
}

// Delete deletes a namespace
func (r GORMNamespaceRepository) Delete(id string) (*domain.Namespace, error) {
	var namespace domain.Namespace
	queryResult := r.Database.Find(&namespace, domain.Namespace{
		ID: id,
	}).Limit(1)
	if queryResult.Error != nil {
		return nil, fmt.Errorf("error finding namespace: %w", queryResult.Error)
	}
	if queryResult.RowsAffected == 0 {
		return nil, fmt.Errorf("namespace with ID %s not found", id)
	}
	namespace.DeletedAt = time.Now()
	deleteResult := r.Database.Save(&namespace)
	if deleteResult.Error != nil {
		return nil, fmt.Errorf("error deleting namespace: %w", deleteResult.Error)
	}
	return &namespace, nil
}
