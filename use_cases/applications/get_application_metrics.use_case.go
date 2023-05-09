package applications

import (
	"cloud-app-hive/domain"
	"fmt"
)

type GetApplicationMetricsUseCase struct {
	ContainerManagerRepository domain.ContainerManagerRepository
}

func (getApplicationMetricsUseCase GetApplicationMetricsUseCase) Execute(application domain.GetApplicationMetrics) ([]domain.ApplicationMetrics, error) {
	metrics, err := getApplicationMetricsUseCase.ContainerManagerRepository.GetApplicationMetrics(application)
	if err != nil {
		return nil, err
	}
	fmt.Println("Metrics: ", metrics)

	return metrics, nil
}
