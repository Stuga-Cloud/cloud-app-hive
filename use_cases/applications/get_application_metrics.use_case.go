package applications

import (
	"cloud-app-hive/domain"
	"fmt"
)

type GetApplicationMetricsUseCase struct {
	ContainerManagerRepository domain.ContainerManagerRepository
}

func (getApplicationMetricsUseCase GetApplicationMetricsUseCase) Execute(appName, appNamespace string) ([]domain.ApplicationMetrics, error) {
	metrics, err := getApplicationMetricsUseCase.ContainerManagerRepository.GetMetricsOfApplication(appNamespace, appName)
	if err != nil {
		return nil, err
	}
	fmt.Println("Metrics: ", metrics)

	return metrics, nil
}
