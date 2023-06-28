package domain

import (
	"cloud-app-hive/controllers/errors"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
)

// ApplicationSecret is a struct that represents a secret
type ApplicationSecret struct {
	Name string `json:"name" validate:"required"`
	Val  string `json:"value" validate:"required"`
}

// ApplicationSecrets is a slice of ApplicationSecret
// swagger:model ApplicationSecrets
type ApplicationSecrets []ApplicationSecret

func (applicationSecrets *ApplicationSecrets) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, &applicationSecrets)
}

func (applicationSecrets *ApplicationSecrets) Value() (driver.Value, error) {
	return json.Marshal(applicationSecrets)
}

const secretNameRegex = "^[a-zA-Z_][a-zA-Z0-9_]*$"

func IsAValidSecretName(name string) bool {
	match, err := regexp.MatchString(secretNameRegex, name)
	if err != nil {
		return false
	}
	return match
}

func (applicationSecrets *ApplicationSecrets) Validate() error {
	for _, secret := range *applicationSecrets {
		if secret.Name == "" {
			return errors.NewInvalidApplicationSecretsError("Name must not be empty")
		}
		if !IsAValidSecretName(secret.Name) {
			return errors.NewInvalidApplicationSecretsError("Name must not contain special characters, it must match the following regex: " + secretNameRegex)
		}
	}
	return nil
}
