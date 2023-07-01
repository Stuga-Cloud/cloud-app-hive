package domain

import (
	"cloud-app-hive/utils"
	"fmt"
)

type ApplicationMetrics struct {
	PodName                 string  `json:"podName"`
	Name                    string  `json:"name"`
	CPUUsage                string  `json:"cpuUsage"`
	MaxCPUUsage             string  `json:"maxCpuUsage"`
	MemoryUsage             string  `json:"memoryUsage"`
	MaxMemoryUsage          string  `json:"maxMemoryUsage"`
	EphemeralStorageUsage   string  `json:"ephemeralStorageUsage"`
	MaxEphemeralStorage     string  `json:"maxEphemeralStorage"`
	PodsUsage               string  `json:"pods"`
	CPUUsageInPercentage    float64 `json:"cpuUsageInPercentage"`
	MemoryUsageInPercentage float64 `json:"memoryUsageInPercentage"`
}

func (applicationMetrics ApplicationMetrics) WithRealLifeReadableUnits() ApplicationMetrics {
	applicationMetrics.CPUUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.CPUUsage)
	applicationMetrics.MaxCPUUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.MaxCPUUsage)
	applicationMetrics.MemoryUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.MemoryUsage)
	applicationMetrics.MaxMemoryUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.MaxMemoryUsage)
	applicationMetrics.EphemeralStorageUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.EphemeralStorageUsage)
	applicationMetrics.MaxEphemeralStorage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.MaxEphemeralStorage)
	applicationMetrics.PodsUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.PodsUsage)
	return applicationMetrics
}

// ToString returns a string representation of the application metrics
func (applicationMetrics ApplicationMetrics) String() string {
	return fmt.Sprintf(
		"PodName: %s, Name: %s, CPUUsage: %s, MaxCPUUsage: %s, MemoryUsage: %s, MaxMemoryUsage: %s, EphemeralStorageUsage: %s, MaxEphemeralStorage: %s, PodsUsage: %s",
		applicationMetrics.PodName,
		applicationMetrics.Name,
		applicationMetrics.CPUUsage,
		applicationMetrics.MaxCPUUsage,
		applicationMetrics.MemoryUsage,
		applicationMetrics.MaxMemoryUsage,
		applicationMetrics.EphemeralStorageUsage,
		applicationMetrics.MaxEphemeralStorage,
		applicationMetrics.PodsUsage,
	)
}
