package requests

import (
	"cloud-app-hive/controllers/validators"
	"github.com/go-playground/validator"
)

// CreateNamespaceRequest is a struct that represents the request body for creating a namespace
type CreateNamespaceRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100" validate:"IsACustomStringForSubdomainValidation"`
	Description string `json:"description" binding:"required,min=3,max=1000"`
	UserID      string `json:"user_id" binding:"required"`
}

func ValidateCreateNamespaceRequest(createNamespaceRequest CreateNamespaceRequest) error {
	validate := validator.New()
	err := validate.RegisterValidation("IsACustomStringForSubdomainValidation", validators.IsACustomStringForSubdomainValidation)
	if err != nil {
		return err
	}

	err = validate.Struct(createNamespaceRequest)
	if err != nil {
		return err
	}
	return nil
}
