package repositories

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GORMApplicationRepository struct {
	Database *gorm.DB
}

// FindApplications returns a list of applications
func (r GORMApplicationRepository) FindApplications(findApplications commands.FindApplications) ([]domain.Application, error) {
	var applications []domain.Application
	query := r.Database
	if findApplications.NamespaceID != nil {
		query = query.Where("namespace_id = ?", findApplications.NamespaceID)
	}
	if findApplications.Name != nil {
		query = query.Where("name = ?", findApplications.Name)
	}
	if findApplications.Image != nil {
		query = query.Where("image = ?", findApplications.Image)
	}
	if findApplications.ApplicationType != nil {
		query = query.Where("application_type = ?", findApplications.ApplicationType)
	}
	if findApplications.IsAutoScaled != nil {
		query = query.Where("scalability_specifications ->> 'is_auto_scaled' = ?", "true")
	}

	result := query.Limit(int(findApplications.Limit)).Offset(int((findApplications.Page - 1) * findApplications.Limit)).Find(&applications)
	if result.Error != nil {
		return nil, fmt.Errorf("error while getting applications: %w", query.Error)
	}
	return applications, nil
}

// FindByID returns an application by its ID
func (r GORMApplicationRepository) FindByID(id string) (*domain.Application, error) {
	app := domain.Application{}
	result := r.Database.Preload("Namespace").Preload("Namespace.Memberships").Limit(1).Find(&app, domain.Application{
		ID: id,
	})

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("application not found with ID %s", id)
	}

	var containerSpecificationsJSON string
	var scalabilitySpecificationsJSON string
	var environmentVariablesJSON string
	var secretsJSON string

	r.Database.Table("applications").Where("id = ?", id).Limit(1).Pluck("container_specifications", &containerSpecificationsJSON)
	r.Database.Table("applications").Where("id = ?", id).Limit(1).Pluck("scalability_specifications", &scalabilitySpecificationsJSON)
	r.Database.Table("applications").Where("id = ?", id).Limit(1).Pluck("environment_variables", &environmentVariablesJSON)
	r.Database.Table("applications").Where("id = ?", id).Limit(1).Pluck("secrets", &secretsJSON)

	var containerSpecifications *domain.ApplicationContainerSpecifications
	var scalabilitySpecifications *domain.ApplicationScalabilitySpecifications
	var environmentVariables *domain.ApplicationEnvironmentVariables
	var secrets *domain.ApplicationSecrets

	if containerSpecificationsJSON != "" && containerSpecificationsJSON != "null" {
		err := json.Unmarshal([]byte(containerSpecificationsJSON), &containerSpecifications)
		if err != nil {
			return nil, err
		}
	}
	if scalabilitySpecificationsJSON != "" && scalabilitySpecificationsJSON != "null" {
		err := json.Unmarshal([]byte(scalabilitySpecificationsJSON), &scalabilitySpecifications)
		if err != nil {
			return nil, err
		}
	}
	if environmentVariablesJSON != "" && environmentVariablesJSON != "null" {
		err := json.Unmarshal([]byte(environmentVariablesJSON), &environmentVariables)
		if err != nil {
			return nil, err
		}
	}
	if secretsJSON != "" && secretsJSON != "null" {
		err := json.Unmarshal([]byte(secretsJSON), &secrets)
		if err != nil {
			return nil, err
		}
	}

	app.ContainerSpecifications = containerSpecifications
	app.ScalabilitySpecifications = scalabilitySpecifications
	app.EnvironmentVariables = environmentVariables
	app.Secrets = secrets

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
	result := r.Database.Where("deleted_at IS NULL").Find(&applications, domain.Application{
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
	app.AdministratorEmail = application.AdministratorEmail

	saveResult := r.Database.Save(&app)
	if saveResult.Error != nil {
		return nil, saveResult.Error
	}
	return &app, nil
}

// Delete deletes an application
func (r GORMApplicationRepository) Delete(id string) (*domain.Application, error) {
	app := domain.Application{}
	foundResult := r.Database.Find(&app, domain.Application{
		ID: id,
	}).Limit(1)
	if foundResult.Error != nil {
		return nil, fmt.Errorf("error while finding application while deleting: %w", foundResult.Error)
	}
	if foundResult.RowsAffected == 0 {
		return nil, fmt.Errorf("application with ID %s not found while deleting", id)
	}

	result := r.Database.Delete(&app)
	if result.Error != nil {
		return nil, fmt.Errorf("error deleting application: %w", result.Error)
	}
	return &app, nil
}

// FindManualScalingApplications returns all applications that are manually scaled
func (r GORMApplicationRepository) FindManualScalingApplications() ([]domain.Application, error) {
	var applications []domain.Application
	result := r.Database.Preload(
		"Namespace",
	).Where(
		"deleted_at IS NULL",
	).Find(&applications, domain.Application{
		ApplicationType: domain.LoadBalanced,
	}).Where(
		"scalabilitySpecifications->>'isAutoScaled' = ?", "false",
	)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding manual scaling applications: %w", result.Error)
	}
	return applications, nil
}
