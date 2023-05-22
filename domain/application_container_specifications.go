package domain

import (
	"cloud-app-hive/domain/errors"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// LimitUnit is an enum that represents the unit of a limit (e.g. KB, MB, GB, etc.)
type LimitUnit string

const (
	// KB is a limit unit that represents a limit in kilobytes
	KB LimitUnit = "KB"
	// MB is a limit unit that represents a limit in megabytes
	MB LimitUnit = "MB"
	// GB is a limit unit that represents a limit in gigabytes
	GB LimitUnit = "GB"
	// TB is a limit unit that represents a limit in terabytes
	TB LimitUnit = "TB"
)

// Scan converts the database value to the custom type
func (limitUnit *LimitUnit) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("failed to scan LimitUnit: value is nil")
	}

	// Assuming the value from the database is stored as a string
	if stringValue, ok := value.(string); ok {
		*limitUnit = LimitUnit(strings.ToUpper(stringValue))
		return nil
	}

	return fmt.Errorf("failed to scan LimitUnit: invalid value type")
}

// Value converts the custom type to a database value
func (limitUnit *LimitUnit) Value() (driver.Value, error) {
	return string(*limitUnit), nil
}

// ContainerLimit is a struct that represents a limit of a container (e.g. CPU, Memory, Storage, etc.)
type ContainerLimit struct {
	Val  int       `json:"value" gorm:"not null"`
	Unit LimitUnit `json:"unit" binding:"omitempty,oneof=KB MB GB TB" gorm:"not null"`
}

func (containerLimit *ContainerLimit) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, &containerLimit)
}

func (containerLimit *ContainerLimit) Value() (driver.Value, error) {
	return json.Marshal(containerLimit)
}

// ApplicationContainerSpecifications is a struct that represents the container characteristics of an application
type ApplicationContainerSpecifications struct {
	CpuLimit     ContainerLimit `json:"cpu_limit" gorm:"json"`
	MemoryLimit  ContainerLimit `json:"memory_limit" gorm:"json"`
	StorageLimit ContainerLimit `json:"storage_limit" gorm:"json"`
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, &applicationContainerSpecifications)
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) Value() (driver.Value, error) {
	return json.Marshal(applicationContainerSpecifications)
}

func (applicationContainerSpecifications ApplicationContainerSpecifications) Validate() error {
	if applicationContainerSpecifications.CpuLimit != (ContainerLimit{}) && applicationContainerSpecifications.CpuLimit.Val <= 0 {
		return errors.NewInvalidApplicationContainerSpecificationsError("CpuLimit must be greater than 0")
	}
	if applicationContainerSpecifications.MemoryLimit != (ContainerLimit{}) && applicationContainerSpecifications.MemoryLimit.Val <= 0 {
		return errors.NewInvalidApplicationContainerSpecificationsError("MemoryLimit must be greater than 0")
	}
	if applicationContainerSpecifications.StorageLimit != (ContainerLimit{}) && applicationContainerSpecifications.StorageLimit.Val <= 0 {
		return errors.NewInvalidApplicationContainerSpecificationsError("StorageLimit must be greater than 0")
	}
	return nil
}
