package dto

import (
	"cloud-app-hive/domain"
	"github.com/go-playground/validator"
)

// UpdateApplicationDto is a struct that represents the request body for creating an application
type UpdateApplicationDto struct {
	Name                      string                                      `json:"name" binding:"required,min=3,max=50" validate:"IsACustomStringForSubdomainValidation"`
	Image                     string                                      `json:"image" binding:"required"`
	Namespace                 string                                      `json:"namespace" binding:"required,min=3,max=50" validate:"IsACustomStringForSubdomainValidation"`
	Port                      int                                         `json:"port" binding:"required, type=number, min=1, max=65535"`
	ApplicationType           domain.ApplicationType                      `json:"application_type" binding:"required,oneof=0 1"`
	EnvironmentVariables      []domain.ApplicationEnvironmentVariable     `json:"environment_variables"`
	Secrets                   []domain.ApplicationSecret                  `json:"secrets"`
	ContainerSpecifications   domain.ApplicationContainerSpecifications   `json:"container_specifications"`
	ScalabilitySpecifications domain.ApplicationScalabilitySpecifications `json:"scalability_specifications"`
}

func ValidateUpdateApplicationDto(updateApplicationDto UpdateApplicationDto) error {
	validate := validator.New()
	err := validate.RegisterValidation("IsACustomStringForSubdomainValidation", IsACustomStringForSubdomainValidation)
	if err != nil {
		return err
	}

	err = validate.Struct(updateApplicationDto)
	if err != nil {
		return err
	}
	return nil
}
