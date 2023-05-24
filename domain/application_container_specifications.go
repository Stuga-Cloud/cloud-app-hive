package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	customErrors "cloud-app-hive/controllers/errors"
)

// LimitUnit is an enum that represents the unit of a limit (e.g. KB, MB, GB, etc.)
type LimitUnit string

const (
	// KB is a limit unit that represents a limit in kilobytes.
	KB LimitUnit = "KB"
	// MB is a limit unit that represents a limit in megabytes.
	MB LimitUnit = "MB"
	// GB is a limit unit that represents a limit in gigabytes.
	GB LimitUnit = "GB"
	// TB is a limit unit that represents a limit in terabytes.
	TB LimitUnit = "TB"
)

// Scan converts the database value to the custom type
func (limitUnit *LimitUnit) Scan(value interface{}) error {
	if value == nil {
		return customErrors.NewLimitUnitScanError("failed to scan LimitUnit: value is nil")
	}

	// Assuming the value from the database is stored as a string
	if stringValue, ok := value.(string); ok {
		*limitUnit = LimitUnit(strings.ToUpper(stringValue))

		return nil
	}

	return customErrors.NewLimitUnitScanError(fmt.Sprintf("failed to scan LimitUnit: value is not a string: %v", value))
}

// Value converts the custom type to a database value.
func (limitUnit *LimitUnit) Value() (driver.Value, error) {
	return string(*limitUnit), nil
}

// ContainerLimit is a struct that represents a limit of a container (e.g. CPU, Memory, Storage, etc.).
type ContainerLimit struct {
	Val  int       `json:"value" gorm:"not null"`
	Unit LimitUnit `json:"unit" binding:"omitempty,oneof=KB MB GB TB" gorm:"not null"`
}

func (containerLimit *ContainerLimit) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return customErrors.NewContainerLimitScanError(fmt.Sprintf("failed to unmarshal JSONB value: %v", value))
	}

	err := json.Unmarshal(bytes, &containerLimit)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", err)
	}

	return nil
}

func (containerLimit *ContainerLimit) Value() (driver.Value, error) {
	res, err := json.Marshal(containerLimit)
	if err != nil {
		// return nil, fmt.Errorf("failed to marshal JSONB value: %v", err)
		return nil, customErrors.NewContainerLimitValueError(fmt.Sprintf("failed to marshal JSONB value: %v", err))
	}

	return res, nil
}

// ApplicationContainerSpecifications is a struct that represents the container characteristics of an application
// swagger:model ApplicationContainerSpecifications
type ApplicationContainerSpecifications struct {
	CPULimit     ContainerLimit `json:"cpuLimit" gorm:"json"`
	MemoryLimit  ContainerLimit `json:"memoryLimit" gorm:"json"`
	StorageLimit ContainerLimit `json:"storageLimit" gorm:"json"`
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return customErrors.NewApplicationContainerSpecificationsScanError(fmt.Sprintf("failed to unmarshal JSONB value: %v", value))
	}

	err := json.Unmarshal(bytes, &applicationContainerSpecifications)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", err)
	}

	return nil
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) Value() (driver.Value, error) {
	res, err := json.Marshal(applicationContainerSpecifications)
	if err != nil {
		return nil, customErrors.NewApplicationContainerSpecificationsValueError(fmt.Sprintf("failed to marshal JSONB value: %v", err))
	}

	return res, nil
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) Validate() error {
	if applicationContainerSpecifications.CPULimit != (ContainerLimit{}) && applicationContainerSpecifications.CPULimit.Val <= 0 {
		return customErrors.NewInvalidApplicationContainerSpecificationsError("CpuLimit must be greater than 0")
	}
	if applicationContainerSpecifications.MemoryLimit != (ContainerLimit{}) && applicationContainerSpecifications.MemoryLimit.Val <= 0 {
		return customErrors.NewInvalidApplicationContainerSpecificationsError("MemoryLimit must be greater than 0")
	}
	if applicationContainerSpecifications.StorageLimit != (ContainerLimit{}) && applicationContainerSpecifications.StorageLimit.Val <= 0 {
		return customErrors.NewInvalidApplicationContainerSpecificationsError("StorageLimit must be greater than 0")
	}
	return nil
}
