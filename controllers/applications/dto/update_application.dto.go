package dto

import (
	"cloud-app-hive/controllers/validators"
	"cloud-app-hive/domain"
	"github.com/go-playground/validator"
)

// UpdateApplicationDto is a struct that represents the request body for creating an application
type UpdateApplicationDto struct {
	Description               string                                      `json:"description"`
	Image                     string                                      `json:"image" binding:"required"`
	Port                      uint32                                      `json:"port" binding:"required,min=1,max=65535"`
	ApplicationType           domain.ApplicationType                      `json:"application_type" binding:"oneof=SINGLE_INSTANCE LOAD_BALANCED"`
	EnvironmentVariables      domain.ApplicationEnvironmentVariables      `json:"environment_variables"`
	Secrets                   domain.ApplicationSecrets                   `json:"secrets"`
	ContainerSpecifications   domain.ApplicationContainerSpecifications   `json:"container_specifications"`
	ScalabilitySpecifications domain.ApplicationScalabilitySpecifications `json:"scalability_specifications"`
}

func ValidateUpdateApplicationDto(updateApplicationDto UpdateApplicationDto) error {
	validate := validator.New()
	err := validate.RegisterValidation("IsACustomStringForSubdomainValidation", validators.IsACustomStringForSubdomainValidation)
	if err != nil {
		return err
	}

	err = validate.Struct(updateApplicationDto)
	if err != nil {
		return err
	}

	err = updateApplicationDto.ContainerSpecifications.Validate()
	if err != nil {
		return err
	}

	err = updateApplicationDto.ScalabilitySpecifications.Validate()
	if err != nil {
		return err
	}

	return nil
}
