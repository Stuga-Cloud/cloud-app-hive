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
	scalabilitySpecifications domain.ApplicationScalabilitySpecifications,
) (bool, error) {
	subject := fmt.Sprintf("Application scalability recommendation - %s in namespace %s", applicationName, namespace)

	cpuLimit := fmt.Sprintf(
		"%d%s",
		specifications.CPULimit.Val,
		specifications.CPULimit.Unit,
	)
	memoryLimit := fmt.Sprintf(
		"%d%s", specifications.MemoryLimit.Val, specifications.MemoryLimit.Unit,
	)
	cpuConfiguredThreshold := scalabilitySpecifications.CpuUsagePercentageThreshold
	memoryConfiguredThreshold := scalabilitySpecifications.MemoryUsagePercentageThreshold
	cpuActualUsage := result.CPUUsageResult.ActualUsage
	memoryActualUsage := result.MemoryUsageResult.ActualUsage

	textBody, htmlBody := s.GetManualScalingRecommendationBody(
		applicationName,
		namespace,
		cpuConfiguredThreshold,
		memoryConfiguredThreshold,
		cpuActualUsage,
		memoryActualUsage,
		cpuLimit,
		memoryLimit,
	)

	err := s.EmailService.Send(to, subject, textBody, htmlBody)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetManualScalingRecommendationBody returns the body of the email sent to the application administrator with the manual scaling recommendation
func (s *ScalabilityNotificationService) GetManualScalingRecommendationBody(
	applicationName string,
	namespace string,
	cpuConfiguredThreshold float64,
	memoryConfiguredThreshold float64,
	cpuActualUsage float64,
	memoryActualUsage float64,
	cpuLimit string,
	memoryLimit string,
) (string, string) {
	body := "The application '" + applicationName + "' in namespace '" + namespace + "' is recommended to be scaled up.\n\n"
	body += "The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.\n\n"
	body += "And the application is currently using " + fmt.Sprintf("%.2f", cpuActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", memoryActualUsage) + "% of its memory resources.\n"
	// and " + fmt.Sprintf("%.2f", result.EphemeralStorageUsageResult.ActualUsage) + "% of its ephemeral storage resources."
	body += "The application is configured to be limited at " +
		cpuLimit + " of CPU resources and " +
		memoryLimit + " of memory resources"

	// fmt.Sprintf(
	// 	"%d%s", specifications.EphemeralStorageLimit.Val, specifications.EphemeralStorageLimit.Unit,
	// ) + " of ephemeral storage resources."
	// body += "The application is configured to have minimum " + fmt.Sprintf("%d", scalabilitySpecifications.MinimumInstanceCount) + " and maximum " + fmt.Sprintf("%d", scalabilitySpecifications.MaximumInstanceCount) + " instances."

	body += "\n\n"
	body += "Best regards,\n"
	body += "The Stuga Cloud Team"

	htmlBody := "<p>The application '" + applicationName + "' in namespace '" + namespace + "' is recommended to be scaled up.</p><br>"
	htmlBody += "<p>The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.</p>"
	htmlBody += "<p>And the application is currently using " + fmt.Sprintf("%.2f", cpuActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", memoryActualUsage) + "% of its memory resources.</p>"

	htmlBody += "<p>The application is configured to be limited at " +
		cpuLimit + " of CPU resources and " +
		memoryLimit + " of memory resources"

	htmlBody += "\n\n"
	htmlBody += "Best regards,<br>"
	htmlBody += "The Stuga Cloud Team"

	return body, htmlBody
}

func (s *ScalabilityNotificationService) SendApplicationScaledUp(
	to string,
	applicationName string,
	namespace string,
	result domain.CompareActualUsageToAcceptedPercentageResult,
	specifications *domain.ApplicationContainerSpecifications,
	scalabilitySpecifications domain.ApplicationScalabilitySpecifications,
) (bool, error) {
	subject := fmt.Sprintf("Application scaled up - %s in namespace %s", applicationName, namespace)

	cpuLimit := fmt.Sprintf(
		"%d%s",
		specifications.CPULimit.Val,
		specifications.CPULimit.Unit,
	)
	memoryLimit := fmt.Sprintf(
		"%d%s", specifications.MemoryLimit.Val, specifications.MemoryLimit.Unit,
	)
	cpuConfiguredThreshold := scalabilitySpecifications.CpuUsagePercentageThreshold
	memoryConfiguredThreshold := scalabilitySpecifications.MemoryUsagePercentageThreshold
	cpuActualUsage := result.CPUUsageResult.ActualUsage
	memoryActualUsage := result.MemoryUsageResult.ActualUsage

	textBody, htmlBody := s.GetAutoScalingBody(
		applicationName,
		namespace,
		cpuConfiguredThreshold,
		memoryConfiguredThreshold,
		cpuActualUsage,
		memoryActualUsage,
		cpuLimit,
		memoryLimit,
	)

	err := s.EmailService.Send(to, subject, textBody, htmlBody)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *ScalabilityNotificationService) GetAutoScalingBody(
	applicationName string,
	namespace string,
	cpuConfiguredThreshold float64,
	memoryConfiguredThreshold float64,
	cpuActualUsage float64,
	memoryActualUsage float64,
	cpuLimit string,
	memoryLimit string,
) (string, string) {
	body := "The application '" + applicationName + "' in namespace '" + namespace + "' has scaled up.\n\n"
	body += "The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.\n\n"
	body += "And the application is currently using " + fmt.Sprintf("%.2f", cpuActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", memoryActualUsage) + "% of its memory resources.\n"
	// and " + fmt.Sprintf("%.2f", result.EphemeralStorageUsageResult.ActualUsage) + "% of its ephemeral storage resources."
	body += "The application is configured to be limited at " +
		cpuLimit + " of CPU resources and " +
		memoryLimit + " of memory resources"

	// fmt.Sprintf(
	// 	"%d%s", specifications.EphemeralStorageLimit.Val, specifications.EphemeralStorageLimit.Unit,
	// ) + " of ephemeral storage resources."
	// body += "The application is configured to have minimum " + fmt.Sprintf("%d", scalabilitySpecifications.MinimumInstanceCount) + " and maximum " + fmt.Sprintf("%d", scalabilitySpecifications.MaximumInstanceCount) + " instances."

	body += "\n\n"
	body += "Best regards,\n"
	body += "The Stuga Cloud Team"

	htmlBody := "<p>The application '" + applicationName + "' in namespace '" + namespace + "' is scaled up.</p><br>"
	htmlBody += "<p>The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.</p>"
	htmlBody += "<p>And the application is currently using " + fmt.Sprintf("%.2f", cpuActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", memoryActualUsage) + "% of its memory resources.</p>"

	htmlBody += "<p>The application is configured to be limited at " +
		cpuLimit + " of CPU resources and " +
		memoryLimit + " of memory resources"

	htmlBody += "\n\n"
	htmlBody += "Best regards,<br>"
	htmlBody += "The Stuga Cloud Team"

	return body, htmlBody
}
