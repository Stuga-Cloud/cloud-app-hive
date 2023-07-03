package schedulers

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/services"
	"cloud-app-hive/use_cases/applications"
	"fmt"
	"os"
	"strconv"
	"time"
)

type NotifyApplicationScalingRecommendationScheduler struct {
	findManualScalingApplicationsUseCase applications.FindManualScalingApplicationsUseCase
	getApplicationMetricsUseCase         applications.GetApplicationMetricsUseCase
	scalabilityNotificationService       services.ScalabilityNotificationService
}

func (scheduler NotifyApplicationScalingRecommendationScheduler) Launch() {
	fmt.Println("Starting 'NotifyApplicationScalingRecommendationScheduler' scheduler...")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic during NotifyApplicationScalingRecommendationScheduler:", r)
			}
		}()
		repeatInterval := getManualScaleAppSchedulerRepeatInterval()
		ticker := time.NewTicker(time.Duration(repeatInterval) * time.Second)
		for {
			select {
			case <-ticker.C:
				// 1. get all applications that are manual scaling
				foundApplications, err := scheduler.findManualScalingApplicationsUseCase.Execute()
				if err != nil {
					fmt.Println("error when try to get manual scaling applications :", err.Error())
					panic("Error when try to get manual scaling applications during NotifyApplicationScalingRecommendationScheduler : " + err.Error())
				}
				if len(foundApplications) == 0 {
					fmt.Println("No manual scaling applications found")
				}

				// 2. parallelize checking each application through the container manager with goroutines
				routines := len(foundApplications)
				done := make(chan bool, routines)
				for _, application := range foundApplications {
					go func(application domain.Application) {
						metrics, err := scheduler.getApplicationMetricsUseCase.Execute(commands.GetApplicationMetrics{
							Name:      application.Name,
							Namespace: application.Namespace.Name,
						})
						if err != nil {
							fmt.Println("error when try to get manual scaling application metrics during NotifyApplicationScalingRecommendationScheduler :", err.Error())
							done <- true
							return
						}
						if len(metrics) == 0 {
							fmt.Println("No metrics found for manual scaling application", application.Name)
							done <- true
							return
						}

						// for _, metric := range metrics {
						// 	fmt.Println("Metrics of application", application.Name, ":", metric.String())
						// }

						var compareActualUsageToAcceptedPercentageResults []domain.CompareActualUsageToAcceptedPercentageResult
						acceptedCPUUsageThreshold := application.ScalabilitySpecifications.Data().CpuUsagePercentageThreshold
						acceptedMemoryUsageThreshold := application.ScalabilitySpecifications.Data().MemoryUsagePercentageThreshold
						for _, metric := range metrics {
							doesExceedCPUAcceptedPercentage, actualCPUUsage := domain.DoesUsageExceedsLimitAndHowMuchActually(metric.CPUUsage, metric.MaxCPUUsage, acceptedCPUUsageThreshold)
							doesExceedMemoryAcceptedPercentage, actualMemoryUsage := domain.DoesUsageExceedsLimitAndHowMuchActually(metric.MemoryUsage, metric.MaxMemoryUsage, acceptedMemoryUsageThreshold)
							// doesExceedEphemeralStorageAcceptedPercentage, actualEphemeralStorageUsage := utils.DoesUsageExceedsLimitAndHowMuchActually(metric.EphemeralStorageUsage, metric.MaxEphemeralStorage, acceptedUsagePercentage)

							compareActualUsageToAcceptedPercentageResult := domain.CompareActualUsageToAcceptedPercentageResult{
								CPUUsageResult: domain.UsageComparisonResult{
									DoesExceedAcceptedPercentage: doesExceedCPUAcceptedPercentage,
									ActualUsage:                  actualCPUUsage,
								},
								MemoryUsageResult: domain.UsageComparisonResult{
									DoesExceedAcceptedPercentage: doesExceedMemoryAcceptedPercentage,
									ActualUsage:                  actualMemoryUsage,
								},
								// EphemeralStorageUsageResult: domain.UsageComparisonResult{
								// 	DoesExceedAcceptedPercentage: doesExceedEphemeralStorageAcceptedPercentage,
								// 	ActualUsage:                  actualEphemeralStorageUsage,
								// },
								PodName: metric.PodName,
							}
							compareActualUsageToAcceptedPercentageResults = append(compareActualUsageToAcceptedPercentageResults, compareActualUsageToAcceptedPercentageResult)
						}

						for _, compareActualUsageToAcceptedPercentageResult := range compareActualUsageToAcceptedPercentageResults {
							// 3. if the application is using more than 80% of its resources, recommend to scale up
							// 4. if the application is using less than 20% of its resources, recommend to scale down
							// 5. send email notification to the application owner
							if compareActualUsageToAcceptedPercentageResult.OneOfTheUsageExceedsAcceptedPercentage() {
								success, err := scheduler.scalabilityNotificationService.SendApplicationScalabilityRecommandationMail(
									application.AdministratorEmail,
									application.Name,
									application.Namespace.Name,
									compareActualUsageToAcceptedPercentageResult,
									application.ContainerSpecifications,
									application.ScalabilitySpecifications.Data(),
								)
								if err != nil {
									fmt.Println("error when try to send application scalability recommandation mail :", err.Error())
									panic("Error when try to send application scalability recommandation mail during NotifyApplicationScalingRecommendationScheduler : " + err.Error())
								}
								if success {
									fmt.Println("Application", application.Name, "has exceeded accepted percentages, email sent to", application.AdministratorEmail)
								} else {
									fmt.Println("Application", application.Name, "has exceeded accepted percentages, but email not sent to", application.AdministratorEmail)
								}
							} else {
								fmt.Println("Application", application.Name, "has not exceeded accepted percentages")
							}
						}

						done <- true
					}(application)
				}
				for i := 0; i < routines; i++ {
					<-done
				}
			}
		}
	}()
}

func getManualScaleAppSchedulerRepeatInterval() int {
	const defaultRepeatInterval = 1
	var repeatInterval int
	SchedulerRecommendApplicationScalingInSeconds := os.Getenv("SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS")
	if SchedulerRecommendApplicationScalingInSeconds == "" {
		fmt.Println("SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS is not set, using default value")
		repeatInterval = defaultRepeatInterval
	}
	repeatInterval, err := strconv.Atoi(SchedulerRecommendApplicationScalingInSeconds)
	if err != nil {
		fmt.Println("Error when convert SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS to int")
		panic(fmt.Sprintf("Error when convert SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS to int during NotifyApplicationScalingRecommendationScheduler : %s", err.Error()))
	}
	return repeatInterval
}
