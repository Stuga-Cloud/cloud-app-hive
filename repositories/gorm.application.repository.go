package repositories

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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
		if *findApplications.IsAutoScaled {
			query = query.Where("scalability_specifications ->> 'isAutoScaled' = ?", "true")
		} else {
			query = query.Where("scalability_specifications ->> 'isAutoScaled' = ?", "false")
		}
	}

	result := query.Limit(int(findApplications.Limit)).Offset(int((findApplications.Page - 1) * findApplications.Limit)).Find(&applications)
	if result.Error != nil {
		return nil, fmt.Errorf("error while getting applications: %w", query.Error)
	}
	return applications, nil
}

// FindByID returns an application by its ID
func (r GORMApplicationRepository) FindByID(id string) (*domain.Application, error) {
	app := &domain.Application{}
	result := r.Database.Preload("Namespace").Preload("Namespace.Memberships").Limit(1).Find(&app, domain.Application{
		ID: id,
	})

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("application not found with ID %s", id)
	}

	app, err := fillApplicationJSONFields(app, r)
	if err != nil {
		return nil, err
	}

	return app, nil
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
	scalabilitySpecs := datatypes.NewJSONType(createApplication.ScalabilitySpecifications)
	app := domain.Application{
		ID:                      uuid.New().String(),
		Name:                    createApplication.Name,
		Description:             createApplication.Description,
		Image:                   createApplication.Image,
		UserID:                  createApplication.UserID,
		NamespaceID:             createApplication.NamespaceID,
		Port:                    createApplication.Port,
		Zone:                    createApplication.Zone,
		ApplicationType:         createApplication.ApplicationType,
		EnvironmentVariables:    &createApplication.EnvironmentVariables,
		Secrets:                 &createApplication.Secrets,
		ContainerSpecifications: &createApplication.ContainerSpecifications,
		// ScalabilitySpecifications: &createApplication.ScalabilitySpecifications,
		// repositories/gorm.application.repository.go:272:34: cannot use scalabilitySpecifications (variable of type *domain.ApplicationScalabilitySpecifications) as *datatypes.JSONType[domain.ApplicationScalabilitySpecifications] value in assignment
		ScalabilitySpecifications: &scalabilitySpecs,
		AdministratorEmail:        createApplication.AdministratorEmail,
	}
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
	queryResult := r.Database.Preload("Namespace").Limit(1).Find(&app, domain.Application{
		ID: applicationID,
	})
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
	scalabilitySpecs := datatypes.NewJSONType(application.ScalabilitySpecifications)
	app.ScalabilitySpecifications = &scalabilitySpecs
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
	foundResult := r.Database.Limit(1).Find(&app, domain.Application{
		ID: id,
	})
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
		"JSON_EXTRACT(scalability_specifications, '$.isAutoScaled') = false",
	).Find(&applications, domain.Application{
		ApplicationType: domain.LoadBalanced,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error finding manual scaling applications: %w", result.Error)
	}

	applications, err := fillApplicationsJSON(applications, r)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

// FindAutoScalingApplications returns all applications that are auto scaled
func (r GORMApplicationRepository) FindAutoScalingApplications() ([]domain.Application, error) {
	var applications []domain.Application
	result := r.Database.Preload(
		"Namespace",
	).Where(
		"JSON_EXTRACT(scalability_specifications, '$.isAutoScaled') = true",
	).Find(&applications, domain.Application{
		ApplicationType: domain.LoadBalanced,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("error finding auto scaling applications: %w", result.Error)
	}

	applications, err := fillApplicationsJSON(applications, r)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func fillApplicationJSONFields(app *domain.Application, r GORMApplicationRepository) (*domain.Application, error) {
	var containerSpecificationsJSON string
	var scalabilitySpecificationsJSON string
	var environmentVariablesJSON string
	var secretsJSON string

	r.Database.Table("applications").Where("id = ?", app.ID).Limit(1).Pluck("container_specifications", &containerSpecificationsJSON)
	r.Database.Table("applications").Where("id = ?", app.ID).Limit(1).Pluck("scalability_specifications", &scalabilitySpecificationsJSON)
	r.Database.Table("applications").Where("id = ?", app.ID).Limit(1).Pluck("environment_variables", &environmentVariablesJSON)
	r.Database.Table("applications").Where("id = ?", app.ID).Limit(1).Pluck("secrets", &secretsJSON)

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
	scalabilitySpecs := datatypes.NewJSONType(*scalabilitySpecifications)
	app.ScalabilitySpecifications = &scalabilitySpecs
	app.EnvironmentVariables = environmentVariables
	app.Secrets = secrets

	return app, nil
}

func fillApplicationsJSON(apps []domain.Application, r GORMApplicationRepository) ([]domain.Application, error) {
	for i, app := range apps {
		app, err := fillApplicationJSONFields(&app, r)
		if err != nil {
			return nil, err
		}
		apps[i] = *app
	}
	return apps, nil
}

// HorizontalScaleUp scales up an application horizontally
func (r GORMApplicationRepository) HorizontalScaleUp(applicationID string) (*domain.Application, error) {
	app := domain.Application{}
	result := r.Database.Find(&app, domain.Application{
		ID: applicationID,
	}).Limit(1)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding application: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("application not found with ID %s", applicationID)
	}

	newNumberOfReplicas := app.ScalabilitySpecifications.Data().Replicas + 1
	scalabilitySpecs := datatypes.NewJSONType(domain.ApplicationScalabilitySpecifications{
		Replicas:                       newNumberOfReplicas,
		IsAutoScaled:                   app.ScalabilitySpecifications.Data().IsAutoScaled,
		CpuUsagePercentageThreshold:    app.ScalabilitySpecifications.Data().CpuUsagePercentageThreshold,
		MemoryUsagePercentageThreshold: app.ScalabilitySpecifications.Data().MemoryUsagePercentageThreshold,
	})

	scalabilitySpecsJSON, err := json.Marshal(scalabilitySpecs)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling scalability specifications: %w", err)
	}

	result = r.Database.Model(&app).Update("scalability_specifications", string(scalabilitySpecsJSON))
	if result.Error != nil {
		return nil, fmt.Errorf("error while updating application: %w", result.Error)
	}

	return &app, nil
}

// HorizontalScaleDown scales down an application horizontally
func (r GORMApplicationRepository) HorizontalScaleDown(applicationID string) (*domain.Application, error) {
	app := domain.Application{}
	result := r.Database.Find(&app, domain.Application{
		ID: applicationID,
	}).Limit(1)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding application: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("application not found with ID %s", applicationID)
	}

	newNumberOfReplicas := app.ScalabilitySpecifications.Data().Replicas - 1
	scalabilitySpecs := datatypes.NewJSONType(domain.ApplicationScalabilitySpecifications{
		Replicas:                       newNumberOfReplicas,
		IsAutoScaled:                   app.ScalabilitySpecifications.Data().IsAutoScaled,
		CpuUsagePercentageThreshold:    app.ScalabilitySpecifications.Data().CpuUsagePercentageThreshold,
		MemoryUsagePercentageThreshold: app.ScalabilitySpecifications.Data().MemoryUsagePercentageThreshold,
	})

	scalabilitySpecsJSON, err := json.Marshal(scalabilitySpecs)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling scalability specifications: %w", err)
	}

	result = r.Database.Model(&app).Update("scalability_specifications", string(scalabilitySpecsJSON))
	if result.Error != nil {
		return nil, fmt.Errorf("error while updating application: %w", result.Error)
	}

	return &app, nil
}

// // VerticalScaleUp scales up an application vertically
// func (r GORMApplicationRepository) VerticalScaleUp(applicationID string) (*domain.Application, error) {
	
// }

// // VerticalScaleDown scales down an application vertically
// func (r GORMApplicationRepository) VerticalScaleDown(applicationID string) (*domain.Application, error) {
// }
