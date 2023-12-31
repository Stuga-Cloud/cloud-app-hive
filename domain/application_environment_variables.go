package domain

import (
	"cloud-app-hive/controllers/errors"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
)

// ApplicationEnvironmentVariable is a struct that represents an environment variable
type ApplicationEnvironmentVariable struct {
	Name string `json:"name" validate:"required"`
	Val  string `json:"value" validate:"required"`
}

// ApplicationEnvironmentVariables is a slice of ApplicationEnvironmentVariable
// swagger:model ApplicationEnvironmentVariables
type ApplicationEnvironmentVariables []ApplicationEnvironmentVariable

func (applicationEnvVariables *ApplicationEnvironmentVariables) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, &applicationEnvVariables)
}

func (applicationEnvVariables *ApplicationEnvironmentVariables) Value() (driver.Value, error) {
	return json.Marshal(applicationEnvVariables)
}

const environmentVariableNameRegex = "^[a-zA-Z_][a-zA-Z0-9_]*$"

func IsAValidEnvironmentVariableName(name string) bool {
	match, err := regexp.MatchString(environmentVariableNameRegex, name)
	if err != nil {
		return false
	}
	return match
}

func (applicationEnvVariables *ApplicationEnvironmentVariables) Validate() error {
	for _, envVariable := range *applicationEnvVariables {
		if envVariable.Name == "" {
			return errors.NewInvalidApplicationEnvironmentVariablesError("Name must not be empty")
		}
		if !IsAValidEnvironmentVariableName(envVariable.Name) {
			return errors.NewInvalidApplicationEnvironmentVariablesError("Name must not contain special characters, it must match the following regex: " + environmentVariableNameRegex)
		}
	}
	return nil
}
