package domain

import (
	"cloud-app-hive/utils"
)

type ApplicationMetrics struct {
	PodName               string `json:"podName"`
	Name                  string `json:"name"`
	CPUUsage              string `json:"cpuUsage"`
	MaxCPUUsage           string `json:"maxCpuUsage"`
	MemoryUsage           string `json:"memoryUsage"`
	MaxMemoryUsage        string `json:"maxMemoryUsage"`
	EphemeralStorageUsage string `json:"ephemeralStorageUsage"`
	MaxEphemeralStorage   string `json:"maxEphemeralStorage"`
	PodsUsage             string `json:"pods"`
}

func (applicationMetrics ApplicationMetrics) WithRealLifeReadableUnits() ApplicationMetrics {
	applicationMetrics.CPUUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.CPUUsage)
	applicationMetrics.MaxCPUUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.MaxCPUUsage)
	applicationMetrics.MemoryUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.MemoryUsage)
	applicationMetrics.MaxMemoryUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.MaxMemoryUsage)
	applicationMetrics.EphemeralStorageUsage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.EphemeralStorageUsage)
	applicationMetrics.MaxEphemeralStorage = utils.ConvertK8sResourceToReadableHumanValueAndUnit(applicationMetrics.MaxEphemeralStorage)
	return applicationMetrics
}
