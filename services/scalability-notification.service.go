package services

import (
	"cloud-app-hive/domain"
	"fmt"
)

type ScalabilityNotificationService struct {
	EmailService EmailService
}

func NewScalabilityNotificationService(emailService EmailService) *ScalabilityNotificationService {
	return &ScalabilityNotificationService{
		EmailService: emailService,
	}
}

// SendApplicationScalabilityRecommandationMail sends an email to the application administrator with the scalability recommendation
func (s *ScalabilityNotificationService) SendApplicationScalabilityRecommandationMail(
	to string,
	applicationName string,
	namespace string,
	result domain.CompareActualUsageToAcceptedPercentageResult,
	specifications *domain.ApplicationContainerSpecifications,
	scalabilitySpecifications *domain.ApplicationScalabilitySpecifications,
) (bool, error) {
	subject := fmt.Sprintf("Application scalability recommendation - %s in namespace %s", applicationName, namespace)
	body := "The application '" + applicationName + "' in namespace '" + namespace + "' is recommended to be scaled up."
	body += "\n"
	body += "The application is currently using " + fmt.Sprintf("%.2f", result.CPUUsageResult.ActualUsage) + "% of its CPU resources, " + fmt.Sprintf("%.2f", result.MemoryUsageResult.ActualUsage) + "% of its memory resources and " + fmt.Sprintf("%.2f", result.EphemeralStorageUsageResult.ActualUsage) + "% of its ephemeral storage resources."
	// TODO Make this work (it's not working because the specifications are not correctly retrieved from database)
	body += "The application is configured to use " +
		fmt.Sprintf(
			"%d%s",
			specifications.CPULimit.Val,
			specifications.CPULimit.Unit,
		) + " of CPU resources, " +
		fmt.Sprintf(
			"%d%s", specifications.MemoryLimit.Val, specifications.MemoryLimit.Unit,
		) + " of memory resources"
	//fmt.Sprintf(
	//	"%d%s", specifications.EphemeralStorageLimit.Val, specifications.EphemeralStorageLimit.Unit,
	//) + " of ephemeral storage resources."
	//body += "The application is configured to have minimum " + fmt.Sprintf("%d", scalabilitySpecifications.MinimumInstanceCount) + " and maximum " + fmt.Sprintf("%d", scalabilitySpecifications.MaximumInstanceCount) + " instances."
	body += "\n\n"
	body += "Best regards,\n"
	body += "The Stuga Cloud Team"
	err := s.EmailService.Send(to, subject, body)
	if err != nil {
		return false, err
	}
	return true, nil
}
