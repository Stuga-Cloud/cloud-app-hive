package repositories

import (
	"fmt"
	"time"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GORMApplicationRepository struct {
	Database *gorm.DB
}

// FindByID returns an application by its ID
func (r GORMApplicationRepository) FindByID(id string) (*domain.Application, error) {
	app := domain.Application{}
	result := r.Database.Preload("Namespace").Find(&app, domain.Application{
		ID: id,
	}).Limit(1)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("application not found with ID %s", id)
	}
	return &app, nil
}

// FindByUserID returns an application by its user ID
func (r GORMApplicationRepository) FindByUserID(userID string) (*domain.Application, error) {
	app := domain.Application{}
	result := r.Database.Find(&app, domain.Application{
		UserID: userID,
	}).Limit(1)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding application: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("application not found with user ID %s", userID)
	}
	return &app, nil
}

// FindByNamespaceIDAndUserID returns an application by namespace ID and user ID
func (r GORMApplicationRepository) FindByNamespaceIDAndUserID(namespaceID string, userID string) ([]domain.Application, error) {
	var applications []domain.Application
	result := r.Database.Find(&applications, domain.Application{
		NamespaceID: namespaceID,
		UserID:      userID,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error finding applications: %w", result.Error)
	}
	return applications, nil
}

// FindByNamespaceIDAndName returns an application by namespace ID and name
func (r GORMApplicationRepository) FindByNamespaceIDAndName(namespaceID string, name string) (*domain.Application, error) {
	app := domain.Application{}
	result := r.Database.First(&app, domain.Application{
		NamespaceID: namespaceID,
		Name:        name,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error finding application: %w", result.Error)
	}
	return &app, nil
}

// Create creates a new application
func (r GORMApplicationRepository) Create(createApplication commands.CreateApplication) (*domain.Application, error) {
	app := createApplication.ToDomain(uuid.New().String())
	result := r.Database.Create(&app)
	if result.Error != nil {
		return nil, fmt.Errorf("error while creating application: %v", result.Error)
	}
	return &app, nil
}

// Update updates an application
func (r GORMApplicationRepository) Update(applicationID string, application commands.UpdateApplication) (*domain.Application, error) {
	app := domain.Application{}
	// Also retrieve namespace linked to application
	queryResult := r.Database.Preload("Namespace").Find(&app, domain.Application{
		ID: applicationID,
	}).Limit(1)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}
	if queryResult.RowsAffected == 0 {
		return nil, fmt.Errorf("application with ID %s not found", applicationID)
	}
	app.Description = application.Description
	app.Image = application.Image
	app.Port = application.Port
	app.ApplicationType = application.ApplicationType
	app.EnvironmentVariables = &application.EnvironmentVariables
	app.Secrets = &application.Secrets
	app.ContainerSpecifications = &application.ContainerSpecifications
	app.ScalabilitySpecifications = &application.ScalabilitySpecifications

	saveResult := r.Database.Save(&app)
	if saveResult.Error != nil {
		return nil, saveResult.Error
	}
	return &app, nil
}

// Delete deletes an application
func (r GORMApplicationRepository) Delete(id string) (*domain.Application, error) {
	app := domain.Application{}
	queryResult := r.Database.Find(&app, domain.Application{
		ID: id,
	}).Limit(1)
	if queryResult.Error != nil {
		return nil, fmt.Errorf("error finding application: %w", queryResult.Error)
	}
	if queryResult.RowsAffected == 0 {
		return nil, fmt.Errorf("application with ID %s not found", id)
	}
	app.DeletedAt = time.Now()
	deleteResult := r.Database.Save(&app)
	if deleteResult.Error != nil {
		return nil, fmt.Errorf("error deleting application: %w", deleteResult.Error)
	}
	return &app, nil
}