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
	Description               string                                      `json:"description" binding:"omitempty,min=3,max=50"`
	Image                     string                                      `json:"image" binding:"required"`
	Registry                  domain.ImageRegistry                        `json:"registry" binding:"required,oneof=dockerhub pcr"`
	NamespaceID               string                                      `json:"namespaceId" binding:"required"`
	UserID                    string                                      `json:"userId" binding:"required"`
	Port                      uint32                                      `json:"port" binding:"required,min=1,max=65535"`
	Zone                      string                                      `json:"zone"`
	ApplicationType           domain.ApplicationType                      `json:"applicationType" binding:"oneof=SINGLE_INSTANCE LOAD_BALANCED" validate:"required"`
	EnvironmentVariables      domain.ApplicationEnvironmentVariables      `json:"environmentVariables"`
	Secrets                   domain.ApplicationSecrets                   `json:"secrets"`
	ContainerSpecifications   domain.ApplicationContainerSpecifications   `json:"containerSpecifications" binding:"required"`
	ScalabilitySpecifications domain.ApplicationScalabilitySpecifications `json:"scalabilitySpecifications" binding:"required"`
	AdministratorEmail        string                                      `json:"administratorEmail" binding:"required,email"`
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

	createApplicationRequest.ContainerSpecifications.SetDefaultValues()
	err = createApplicationRequest.ScalabilitySpecifications.Validate()
	if err != nil {
		return err
	}

	return nil
}
