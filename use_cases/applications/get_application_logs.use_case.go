package applications

import (
	"cloud-app-hive/domain"
	"fmt"
)

type GetApplicationLogsUseCase struct {
	ContainerManagerRepository domain.ContainerManagerRepository
}

func (getApplicationLogsUseCase GetApplicationLogsUseCase) Execute(application domain.GetApplicationLogs) ([]domain.ApplicationLogs, error) {
	logs, err := getApplicationLogsUseCase.ContainerManagerRepository.GetApplicationLogs(application)
	if err != nil {
		return []domain.ApplicationLogs{}, fmt.Errorf("error while getting logs: %w", err)
	}

	return logs, nil
}
