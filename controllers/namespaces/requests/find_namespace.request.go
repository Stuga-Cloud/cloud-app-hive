package requests

import (
	"cloud-app-hive/controllers/validators"
	"github.com/go-playground/validator"
)

// FindNamespacesRequest is a struct that represents the request body for finding a namespace
type FindNamespacesRequest struct {
	Name    *string `form:"name" validate:"omitempty,IsACustomStringForSubdomainValidation"`
	UserID  string  `form:"userId" validate:"required"`
	Page    int     `form:"page"`
	PerPage int     `form:"per_page"`
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
