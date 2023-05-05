package dto

import (
	"github.com/go-playground/validator"
)

// CreateApplicationDto is a struct that represents the request body for creating an application
type CreateApplicationDto struct {
	Name      string `json:"name" binding:"required,min=3,max=50"`
	Image     string `json:"image" binding:"required"`
	Namespace string `json:"namespace" binding:"required,min=3,max=50"`
	//EnvVars []string `json:"env_vars"`
}

func ValidateCreateApplicationDto(createApplicationDto CreateApplicationDto) error {
	validate := validator.New()
	err := validate.Struct(createApplicationDto)
	if err != nil {
		return err
	}
	return nil
}
