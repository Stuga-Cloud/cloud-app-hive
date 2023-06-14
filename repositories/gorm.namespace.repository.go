package repositories

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GORMNamespaceRepository struct {
	Database *gorm.DB
}

// FindByID returns a namespace by its ID
func (r GORMNamespaceRepository) FindByID(id string) (*domain.Namespace, error) {
	app := domain.Namespace{}
	result := r.Database.Preload("Memberships").Preload("Applications").Find(&app, domain.Namespace{
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
	namespace := domain.Namespace{}
	result := r.Database.Preload("Memberships").Find(&namespace, domain.Namespace{
		Name: name,
	}).Limit(1)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding namespace: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("namespace not found with name %s", name)
	}
	return &namespace, nil
}

// Find returns a list of namespaces
func (r GORMNamespaceRepository) Find(findNamespaces commands.FindNamespaces) ([]domain.Namespace, error) {
	var foundNamespaces []domain.Namespace
	// It should match all the filters

	result := r.Database.Preload("Memberships").Preload("Applications")
	if findNamespaces.Name != nil && *findNamespaces.Name != "" {
		result = result.Where("name = ?", findNamespaces.Name)
	}
	result = result.Find(&foundNamespaces).Limit(findNamespaces.PerPage).Offset(findNamespaces.PerPage * (findNamespaces.Page - 1))
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
	err := r.Database.Transaction(func(tx *gorm.DB) error {
		exists, err := r.ExistsByName(namespace.Name)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("namespace with name %s already exists", namespace.Name)
		}

		result := tx.Create(&namespace)
		if result.Error != nil {
			return fmt.Errorf("error creating namespace: %w", result.Error)
		}

		membership := domain.NamespaceMembership{
			ID:          uuid.New().String(),
			UserID:      namespace.UserID,
			NamespaceID: namespace.ID,
			Role:        domain.RoleAdmin,
		}
		result = tx.Create(&membership)
		if result.Error != nil {
			return fmt.Errorf("error creating namespace membership: %w", result.Error)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var createdNamespace domain.Namespace
	result := r.Database.Preload("Memberships").Find(&createdNamespace, domain.Namespace{
		ID: namespace.ID,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error finding created namespace: %w", result.Error)
	}
	return &createdNamespace, nil
}

// Delete deletes a namespace
func (r GORMNamespaceRepository) Delete(id string, userId string) (*domain.Namespace, error) {
	var namespace domain.Namespace
	queryResult := r.Database.Preload("Memberships").Find(&namespace, domain.Namespace{
		ID: id,
	}).Limit(1)
	if queryResult.Error != nil {
		return nil, fmt.Errorf("error finding namespace: %w", queryResult.Error)
	}
	if queryResult.RowsAffected == 0 {
		return nil, fmt.Errorf("namespace with ID %s not found", id)
	}

	// check if user is admin
	isAdmin := false
	for _, membership := range namespace.Memberships {
		if membership.UserID == userId && membership.Role == domain.RoleAdmin {
			isAdmin = true
		}
	}
	if !isAdmin {
		return nil, fmt.Errorf("user is not admin of namespace %s", id)
	}

	result := r.Database.Delete(&namespace)
	if result.Error != nil {
		return nil, fmt.Errorf("error deleting namespace: %w", result.Error)
	}

	// ALso delete all memberships
	memberships := []domain.NamespaceMembership{}
	result = r.Database.Find(&memberships, domain.NamespaceMembership{
		NamespaceID: id,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error finding memberships: %w", result.Error)
	}
	for _, membership := range memberships {
		result = r.Database.Delete(&membership)
		if result.Error != nil {
			return nil, fmt.Errorf("error deleting membership: %w", result.Error)
		}
	}

	// Delete applications
	applications := []domain.Application{}
	result = r.Database.Find(&applications, domain.Application{
		NamespaceID: id,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error finding applications: %w", result.Error)
	}

	for _, application := range applications {
		result = r.Database.Delete(&application)
		if result.Error != nil {
			return nil, fmt.Errorf("error deleting application: %w", result.Error)
		}
	}

	return &namespace, nil
}

// Update updates a namespace
func (r GORMNamespaceRepository) Update(updateNamespace commands.UpdateNamespace) (*domain.Namespace, error) {
	var namespace domain.Namespace
	queryResult := r.Database.Preload("Memberships").Find(&namespace, domain.Namespace{
		ID: updateNamespace.ID,
	}).Limit(1)
	if queryResult.Error != nil {
		return nil, fmt.Errorf("error finding namespace: %w", queryResult.Error)
	}
	if queryResult.RowsAffected == 0 {
		return nil, fmt.Errorf("namespace with ID %s not found", updateNamespace.ID)
	}

	namespace.Description = updateNamespace.Description
	namespace.UserID = updateNamespace.UserID

	result := r.Database.Save(&namespace)
	if result.Error != nil {
		return nil, fmt.Errorf("error updating namespace: %w", result.Error)
	}
	return &namespace, nil
}
