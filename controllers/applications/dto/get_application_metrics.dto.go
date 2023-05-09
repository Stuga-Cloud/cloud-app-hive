package dto

import "github.com/go-playground/validator"

// GetApplicationMetricsDto is a struct that represents the query params for getting metrics for an application
type GetApplicationMetricsDto struct {
	Name      string `url:"name" validate:"IsACustomStringForSubdomainValidation"`
	Namespace string `url:"namespace" validate:"IsACustomStringForSubdomainValidation"`
}

func ValidateGetApplicationMetricsDto(getApplicationMetricsDto GetApplicationMetricsDto) error {
	validate := validator.New()
	err := validate.RegisterValidation("IsACustomStringForSubdomainValidation", IsACustomStringForSubdomainValidation)
	if err != nil {
		return err
	}

	err = validate.Struct(getApplicationMetricsDto)
	if err != nil {
		return err
	}
	return nil
}
