package dto

import (
	"cloud-app-hive/controllers/validators"
	"github.com/go-playground/validator"
)

// CreateNamespaceDto is a struct that represents the request body for creating a namespace
type CreateNamespaceDto struct {
	Name        string `json:"name" binding:"required,min=3,max=100" validate:"IsACustomStringForSubdomainValidation"`
	Description string `json:"description" binding:"required,min=3,max=1000"`
	UserID      string `json:"user_id" binding:"required"`
}

func ValidateCreateNamespaceDto(createNamespaceDto CreateNamespaceDto) error {
	validate := validator.New()
	err := validate.RegisterValidation("IsACustomStringForSubdomainValidation", validators.IsACustomStringForSubdomainValidation)
	if err != nil {
		return err
	}

	err = validate.Struct(createNamespaceDto)
	if err != nil {
		return err
	}
	return nil
}
