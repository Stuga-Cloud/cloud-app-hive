package use_cases

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/repositories"
)

type GetClusterMetricsUseCase struct {
	ContainerManagerRepository repositories.ContainerManagerRepository
}

func (getClusterMetricsUseCase GetClusterMetricsUseCase) Execute() (*domain.ClusterMetrics, error) {
	clusterState, err := getClusterMetricsUseCase.ContainerManagerRepository.GetClusterMetrics()
	if err != nil {
		return nil, fmt.Errorf("error while applying application: %w", err)
	}
	return clusterState, nil
}
