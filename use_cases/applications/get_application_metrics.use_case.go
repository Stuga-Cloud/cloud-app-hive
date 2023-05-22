package applications

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
	"fmt"
)

type GetApplicationMetricsUseCase struct {
	ContainerManagerRepository repositories.ContainerManagerRepository
}

func (getApplicationMetricsUseCase GetApplicationMetricsUseCase) Execute(application commands.GetApplicationMetrics) ([]domain.ApplicationMetrics, error) {
	metrics, err := getApplicationMetricsUseCase.ContainerManagerRepository.GetApplicationMetrics(application)
	if err != nil {
		return nil, fmt.Errorf("error getting application metrics: %w", err)
	}
	return metrics, nil
}
