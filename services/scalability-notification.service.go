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
	results []domain.CompareActualUsageToAcceptedPercentageResult,
	specifications domain.ApplicationContainerSpecifications,
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

	textBody, htmlBody := s.GetManualScalingRecommendationBody(
		applicationName,
		namespace,
		cpuConfiguredThreshold,
		memoryConfiguredThreshold,
		results,
		cpuLimit,
		memoryLimit,
	)

	err := s.EmailService.Send(to, subject, textBody, htmlBody, []string{})
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
	results []domain.CompareActualUsageToAcceptedPercentageResult,
	cpuLimit string,
	memoryLimit string,
) (string, string) {
	body := "The application '" + applicationName + "' in namespace '" + namespace + "' is recommended to be scaled up.\n\n"
	body += "The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.\n\n"

	body += "And the application is currently using :\n"
	for _, result := range results {
		body += "- " + result.PodName + " : " + fmt.Sprintf("%.2f", result.CPUUsageResult.ActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", result.MemoryUsageResult.ActualUsage) + "% of its memory resources.\n"
	}

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
	htmlBody += "<p>And the application is currently using :</p>"
	for _, result := range results {
		htmlBody += "<p>- " + result.PodName + " : " + fmt.Sprintf("%.2f", result.CPUUsageResult.ActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", result.MemoryUsageResult.ActualUsage) + "% of its memory resources.</p>"
	}

	htmlBody += "<p>The application is configured to be limited at " +
		cpuLimit + " of CPU resources and " +
		memoryLimit + " of memory resources"

	htmlBody += "\n\n"
	htmlBody += "Best regards,<br>"
	htmlBody += "The Stuga Cloud Team"

	return body, htmlBody
}

func (s *ScalabilityNotificationService) GetCannotScaleApplicationVerticallyBody(
	applicationName string,
	namespace string,
	cpuConfiguredThreshold float64,
	memoryConfiguredThreshold float64,
	currentUsageResult []domain.CompareActualUsageToAcceptedPercentageResult,
	cpuLimit string,
	memoryLimit string,
) (string, string) {
	body := "The application '" + applicationName + "' in namespace '" + namespace + "' cannot be scaled up more, it has reached the maximum number of CPU, memory resources and replicas.\n\n"
	body += "The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.\n\n"

	body += "The application is currently using :\n"
	for _, result := range currentUsageResult {
		body += "- " + result.PodName + " : " + fmt.Sprintf("%.2f", result.CPUUsageResult.ActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", result.MemoryUsageResult.ActualUsage) + "% of its memory resources.\n"
	}

	// and " + fmt.Sprintf("%.2f", result.EphemeralStorageUsageResult.ActualUsage) + "% of its ephemeral storage resources."
	body += "The application is configured to be limited at " +
		cpuLimit + " of CPU resources and " +
		memoryLimit + " of memory resources and " +
		fmt.Sprintf("%d", domain.MaxNumberOfReplicas) + " replicas."

	body += "\n\n"
	body += "Best regards,\n"
	body += "The Stuga Cloud Team"

	htmlBody := "<p>The application '" + applicationName + "' in namespace '" + namespace + "' cannot be scaled up more, it has reached the maximum number of CPU and memory resources and replicas.</p><br>"
	htmlBody += "<p>The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.</p>"
	htmlBody += "<p>The application is currently using :</p>"
	for _, result := range currentUsageResult {
		htmlBody += "<p>- " + result.PodName + " : " + fmt.Sprintf("%.2f", result.CPUUsageResult.ActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", result.MemoryUsageResult.ActualUsage) + "% of its memory resources.</p>"
	}

	htmlBody += "<p>The application is configured to be limited at " +
		cpuLimit + " of CPU resources and " +
		memoryLimit + " of memory resources and " +
		fmt.Sprintf("%d", domain.MaxNumberOfReplicas) + " replicas."

	htmlBody += "\n\n"
	htmlBody += "Best regards,<br>"
	htmlBody += "The Stuga Cloud Team"

	return body, htmlBody
}

func (s *ScalabilityNotificationService) SendCannotScaleApplicationVertically(
	to string,
	applicationName string,
	namespace string,
	currentUsageResult []domain.CompareActualUsageToAcceptedPercentageResult,
	specifications domain.ApplicationContainerSpecifications,
	scalabilitySpecifications domain.ApplicationScalabilitySpecifications,
) (bool, error) {
	subject := fmt.Sprintf("Cannot scale application anymore - %s in namespace %s", applicationName, namespace)

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

	textBody, htmlBody := s.GetCannotScaleApplicationVerticallyBody(
		applicationName,
		namespace,
		cpuConfiguredThreshold,
		memoryConfiguredThreshold,
		currentUsageResult,
		cpuLimit,
		memoryLimit,
	)

	err := s.EmailService.Send(to, subject, textBody, htmlBody, []string{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *ScalabilityNotificationService) SendApplicationVerticallyScaledUp(
	to string,
	applicationName string,
	namespace string,
	currentUsageResult []domain.CompareActualUsageToAcceptedPercentageResult,
	specifications domain.ApplicationContainerSpecifications,
	scalabilitySpecifications domain.ApplicationScalabilitySpecifications,
) (bool, error) {
	subject := fmt.Sprintf("Application vertically scaled up - %s in namespace %s", applicationName, namespace)

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

	textBody, htmlBody := s.GetAutoScalingBody(
		applicationName,
		namespace,
		cpuConfiguredThreshold,
		memoryConfiguredThreshold,
		currentUsageResult,
		cpuLimit,
		memoryLimit,
		scalabilitySpecifications.Replicas,
	)

	err := s.EmailService.Send(to, subject, textBody, htmlBody, []string{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *ScalabilityNotificationService) SendApplicationHorizontallyScaledUp(
	to string,
	applicationName string,
	namespace string,
	currentUsageResult []domain.CompareActualUsageToAcceptedPercentageResult,
	specifications domain.ApplicationContainerSpecifications,
	scalabilitySpecifications domain.ApplicationScalabilitySpecifications,
) (bool, error) {
	subject := fmt.Sprintf("Application horizontally scaled up - %s in namespace %s", applicationName, namespace)

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

	textBody, htmlBody := s.GetAutoScalingBody(
		applicationName,
		namespace,
		cpuConfiguredThreshold,
		memoryConfiguredThreshold,
		currentUsageResult,
		cpuLimit,
		memoryLimit,
		scalabilitySpecifications.Replicas,
	)

	err := s.EmailService.Send(to, subject, textBody, htmlBody, []string{})
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
	currentUsageResult []domain.CompareActualUsageToAcceptedPercentageResult,
	cpuLimit string,
	memoryLimit string,
	replicas int32,
) (string, string) {
	body := "The application '" + applicationName + "' in namespace '" + namespace + "' has scaled up.\n\n"
	body += "The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.\n\n"

	body += "The application is currently using :\n"
	for _, result := range currentUsageResult {
		body += "- " + result.PodName + " : " + fmt.Sprintf("%.2f", result.CPUUsageResult.ActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", result.MemoryUsageResult.ActualUsage) + "% of its memory resources.\n"
	}

	// and " + fmt.Sprintf("%.2f", result.EphemeralStorageUsageResult.ActualUsage) + "% of its ephemeral storage resources."
	body += "The application is now configured to be limited at " +
		cpuLimit + " of CPU resources and " +
		memoryLimit + " of memory resources"

	// fmt.Sprintf(
	// 	"%d%s", specifications.EphemeralStorageLimit.Val, specifications.EphemeralStorageLimit.Unit,
	// ) + " of ephemeral storage resources."
	// body += "The application is configured to have minimum " + fmt.Sprintf("%d", scalabilitySpecifications.MinimumInstanceCount) + " and maximum " + fmt.Sprintf("%d", scalabilitySpecifications.MaximumInstanceCount) + " instances."

	body += "\n\n"
	body += "Best regards,\n"
	body += "The Stuga Cloud Team"

	htmlBody := "<p>The application '" + applicationName + "' in namespace '" + namespace + "' has scaled up.</p><br>"
	htmlBody += "<p>The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.</p>"

	htmlBody += "<p>The application is currently using :</p>"
	for _, result := range currentUsageResult {
		htmlBody += "<p>- " + result.PodName + " : " + fmt.Sprintf("%.2f", result.CPUUsageResult.ActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", result.MemoryUsageResult.ActualUsage) + "% of its memory resources.</p>"
	}

	htmlBody += "<p>The application is configured to be limited at " +
		cpuLimit + " of CPU resources and " +
		memoryLimit + " of memory resources"

	htmlBody += "<br><br>"
	htmlBody += "Best regards,<br>"
	htmlBody += "The Stuga Cloud Team"

	return body, htmlBody
}

func (s *ScalabilityNotificationService) GetCannotScaleUpBody(
	applicationName string,
	namespace string,
	cpuConfiguredThreshold float64,
	memoryConfiguredThreshold float64,
	currentUsageResult []domain.CompareActualUsageToAcceptedPercentageResult,
	cpuLimit string,
	memoryLimit string,
) (string, string) {
	body := "The application '" + applicationName + "' in namespace '" + namespace + "' cannot be scaled up more, it has reached the maximum number of replicas (" + fmt.Sprintf("%d", domain.MaxNumberOfReplicas) + "), CPU and memory resources.\n\n"
	body += "The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.\n\n"

	body += "The application is currently using :\n"
	for _, result := range currentUsageResult {
		body += "- " + result.PodName + " : " + fmt.Sprintf("%.2f", result.CPUUsageResult.ActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", result.MemoryUsageResult.ActualUsage) + "% of its memory resources.\n"
	}

	// and " + fmt.Sprintf("%.2f", result.EphemeralStorageUsageResult.ActualUsage) + "% of its ephemeral storage resources."
	body += "The application is configured to be limited at " +
		cpuLimit + " of CPU resources and " +
		memoryLimit + " of memory resources and " +
		fmt.Sprintf("%d", domain.MaxNumberOfReplicas) + " replicas."

	body += "\n\n"
	body += "Best regards,\n"
	body += "The Stuga Cloud Team"

	htmlBody := "<p>The application '" + applicationName + "' in namespace '" + namespace + "' cannot be scaled up more, it has reached the maximum number of replicas (" + fmt.Sprintf("%d", domain.MaxNumberOfReplicas) + ") and CPU and memory resources.</p><br>"
	htmlBody += "<p>The configured thresholds are " + fmt.Sprintf("%.2f", cpuConfiguredThreshold) + "% of CPU resources and " + fmt.Sprintf("%.2f", memoryConfiguredThreshold) + "% of memory resources.</p>"

	htmlBody += "<p>The application is currently using :</p>"
	for _, result := range currentUsageResult {
		htmlBody += "<p>- " + result.PodName + " : " + fmt.Sprintf("%.2f", result.CPUUsageResult.ActualUsage) + "% of its CPU resources and " + fmt.Sprintf("%.2f", result.MemoryUsageResult.ActualUsage) + "% of its memory resources.</p>"
	}

	htmlBody += "<p>The application is configured to be limited at " +
		cpuLimit +
		" of CPU resources and " +
		memoryLimit +
		" of memory resources and " +
		fmt.Sprintf("%d", domain.MaxNumberOfReplicas) +
		" replicas."

	htmlBody += "<br><br>"
	htmlBody += "Best regards,<br>"
	htmlBody += "The Stuga Cloud Team"

	return body, htmlBody
}

func (s *ScalabilityNotificationService) SendApplicationCannotBeScaledUp(
	to string,
	applicationName string,
	namespace string,
	currentUsageResults []domain.CompareActualUsageToAcceptedPercentageResult,
	specifications domain.ApplicationContainerSpecifications,
	scalabilitySpecifications domain.ApplicationScalabilitySpecifications,
) (bool, error) {
	subject := fmt.Sprintf("Application cannot be scaled up more - %s in namespace %s", applicationName, namespace)

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

	textBody, htmlBody := s.GetCannotScaleUpBody(
		applicationName,
		namespace,
		cpuConfiguredThreshold,
		memoryConfiguredThreshold,
		currentUsageResults,
		cpuLimit,
		memoryLimit,
	)

	err := s.EmailService.Send(to, subject, textBody, htmlBody, []string{})
	if err != nil {
		return false, err
	}
	return true, nil
}
