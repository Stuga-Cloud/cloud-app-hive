package applications

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type GetApplicationLogsUseCase struct {
	ContainerManagerRepository repositories.ContainerManagerRepository
}

func (getApplicationLogsUseCase GetApplicationLogsUseCase) Execute(application commands.GetApplicationLogs) ([]domain.ApplicationLogs, error) {
	logs, err := getApplicationLogsUseCase.ContainerManagerRepository.GetApplicationLogs(application)
	if err != nil {
		return []domain.ApplicationLogs{}, fmt.Errorf("error while getting logs: %w", err)
	}

	return logs, nil
}
