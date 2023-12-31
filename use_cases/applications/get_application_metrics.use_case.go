package applications

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type GetApplicationMetricsUseCase struct {
	ContainerManagerRepository repositories.ContainerManagerRepository
}

func (getApplicationMetricsUseCase GetApplicationMetricsUseCase) Execute(application commands.GetApplicationMetrics) ([]domain.ApplicationMetrics, error) {
	metrics, err := getApplicationMetricsUseCase.ContainerManagerRepository.GetApplicationMetrics(application)
	if err != nil {
		return nil, fmt.Errorf("error getting application metrics: %w", err)
	}
	readableMetrics := make([]domain.ApplicationMetrics, len(metrics))
	for i, podMetrics := range metrics {
		readableMetrics[i] = podMetrics.WithRealLifeReadableUnits()
		// fmt.Printf("CPU usage: %s and CPU limit: %s\n", readableMetrics[i].CPUUsage, readableMetrics[i].MaxCPUUsage)
		_, cpuUsageInPercentage := domain.DoesUsageExceedsLimitAndHowMuchActually(readableMetrics[i].CPUUsage, readableMetrics[i].MaxCPUUsage, 0)
		_, memoryUsageInPercentage := domain.DoesUsageExceedsLimitAndHowMuchActually(readableMetrics[i].MemoryUsage, readableMetrics[i].MaxMemoryUsage, 0)
		readableMetrics[i].CPUUsageInPercentage = cpuUsageInPercentage
		readableMetrics[i].MemoryUsageInPercentage = memoryUsageInPercentage
	}
	return readableMetrics, nil
}
