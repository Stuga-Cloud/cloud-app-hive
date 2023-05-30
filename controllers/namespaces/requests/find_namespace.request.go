package requests

import (
	"cloud-app-hive/controllers/validators"
	"github.com/go-playground/validator"
)

// FindNamespacesRequest is a struct that represents the request body for finding a namespace
type FindNamespacesRequest struct {
	Name    string `json:"name" validate:"IsACustomStringForSubdomainValidation"`
	UserID  string `json:"userId"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

func ValidateFindNamespacesRequest(findNamespacesRequest FindNamespacesRequest) error {
	validate := validator.New()
	err := validate.RegisterValidation("IsACustomStringForSubdomainValidation", validators.IsACustomStringForSubdomainValidation)
	if err != nil {
		return err
	}

	err = validate.Struct(findNamespacesRequest)
	if err != nil {
		return err
	}
	return nil
}
