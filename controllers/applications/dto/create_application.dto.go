package dto

import (
	"cloud-app-hive/controllers/validators"
	"cloud-app-hive/domain"
	"github.com/go-playground/validator"
)

// CreateApplicationDto is a struct that represents the request body for creating an application
type CreateApplicationDto struct {
	Name                      string                                      `json:"name" binding:"required,min=3,max=50" validate:"IsACustomStringForSubdomainValidation"`
	Image                     string                                      `json:"image" binding:"required"`
	NamespaceID               string                                      `json:"namespace_id" binding:"required"`
	UserID                    string                                      `json:"user_id" binding:"required"`
	Port                      uint32                                      `json:"port" binding:"required,min=1,max=65535"`
	ApplicationType           domain.ApplicationType                      `json:"application_type" binding:"oneof=SINGLE_INSTANCE LOAD_BALANCED"`
	EnvironmentVariables      domain.ApplicationEnvironmentVariables      `json:"environment_variables"`
	Secrets                   domain.ApplicationSecrets                   `json:"secrets"`
	ContainerSpecifications   domain.ApplicationContainerSpecifications   `json:"container_specifications"`
	ScalabilitySpecifications domain.ApplicationScalabilitySpecifications `json:"scalability_specifications"`
}

func ValidateCreateApplicationDto(createApplicationDto CreateApplicationDto) error {
	validate := validator.New()
	err := validate.RegisterValidation("IsACustomStringForSubdomainValidation", validators.IsACustomStringForSubdomainValidation)
	if err != nil {
		return err
	}

	err = validate.Struct(createApplicationDto)
	if err != nil {
		return err
	}

	err = createApplicationDto.ContainerSpecifications.Validate()
	if err != nil {
		return err
	}

	err = createApplicationDto.ScalabilitySpecifications.Validate()
	if err != nil {
		return err
	}

	return nil
}
