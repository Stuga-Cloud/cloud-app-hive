package requests

import (
	"cloud-app-hive/controllers/validators"
	"cloud-app-hive/domain"
	"github.com/go-playground/validator"
)

// CreateApplicationRequest is a struct that represents the request body for creating an application
// swagger:model CreateApplicationRequest
type CreateApplicationRequest struct {
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

func ValidateCreateApplicationRequest(createApplicationRequest CreateApplicationRequest) error {
	validate := validator.New()
	err := validate.RegisterValidation("IsACustomStringForSubdomainValidation", validators.IsACustomStringForSubdomainValidation)
	if err != nil {
		return err
	}

	err = validate.Struct(createApplicationRequest)
	if err != nil {
		return err
	}

	err = createApplicationRequest.ContainerSpecifications.Validate()
	if err != nil {
		return err
	}

	err = createApplicationRequest.ScalabilitySpecifications.Validate()
	if err != nil {
		return err
	}

	return nil
}
