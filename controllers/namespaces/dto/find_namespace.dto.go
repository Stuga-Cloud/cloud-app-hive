package dto

import (
	"cloud-app-hive/controllers/validators"
	"github.com/go-playground/validator"
)

// FindNamespacesDto is a struct that represents the request body for finding a namespace
type FindNamespacesDto struct {
	Name    string `json:"name" validate:"IsACustomStringForSubdomainValidation"`
	UserID  string `json:"user_id"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

func ValidateFindNamespacesDto(findNamespacesDto FindNamespacesDto) error {
	validate := validator.New()
	err := validate.RegisterValidation("IsACustomStringForSubdomainValidation", validators.IsACustomStringForSubdomainValidation)
	if err != nil {
		return err
	}

	err = validate.Struct(findNamespacesDto)
	if err != nil {
		return err
	}
	return nil
}
