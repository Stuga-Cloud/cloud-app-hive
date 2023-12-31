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

type NotifyApplicationManualScalingRecommendationScheduler struct {
	findManualScalingApplicationsUseCase applications.FindManualScalingApplicationsUseCase
	getApplicationMetricsUseCase         applications.GetApplicationMetricsUseCase
	scalabilityNotificationService       services.ScalabilityNotificationService
}

func (scheduler NotifyApplicationManualScalingRecommendationScheduler) Launch() {
	fmt.Println("Starting 'NotifyApplicationManualScalingRecommendationScheduler' scheduler...")
	go func() {
		repeatInterval, err := getManualScaleAppSchedulerRepeatInterval()
		if err != nil {
			fmt.Println("Error when try to get notify application manual scaling recommendation scheduler repeat interval :", err.Error())
			return
		}
		ticker := time.NewTicker(time.Duration(repeatInterval) * time.Second)

		lastNotifiedDatetimeByApplication := make(map[string]time.Time)

		for {
			select {
			case <-ticker.C:
				// 1. get all applications that are manual scaling
				foundApplications, err := scheduler.findManualScalingApplicationsUseCase.Execute()
				if err != nil {
					fmt.Println("error when try to get manual scaling applications during NotifyApplicationManualScalingRecommendationScheduler :", err.Error())
					continue
				}
				if len(foundApplications) == 0 {
					// fmt.Println("No manual scaling applications found")
					continue
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
							fmt.Println("error when try to get manual scaling application metrics during NotifyApplicationManualScalingRecommendationScheduler :", err.Error())
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

						oneOfTheUsageExceedsAcceptedPercentage := false
						for _, compareActualUsageToAcceptedPercentageResult := range compareActualUsageToAcceptedPercentageResults {
							// 3. if the application is using more than 80% of its resources, recommend to scale up
							// 4. if the application is using less than 20% of its resources, recommend to scale down
							// 5. send email notification to the application owner
							if compareActualUsageToAcceptedPercentageResult.OneOfTheUsageExceedsAcceptedPercentage() {
								oneOfTheUsageExceedsAcceptedPercentage = true
								break
							}
						}

						if oneOfTheUsageExceedsAcceptedPercentage {
							// If since only 4 hours
							if time.Since(lastNotifiedDatetimeByApplication[application.ID]).Hours() < 4 {
								fmt.Println("Manual application", application.Name, "has exceeded accepted percentages, but last notification was less than 4 hours ago")
								done <- true
								return
							}

							// jsonApplication, err := json.Marshal(application)
							//
							//	if err != nil {
							//		fmt.Println("error when try to marshal application to json during NotifyApplicationManualScalingRecommendationScheduler :", err.Error())
							//		done <- true
							//		return
							//	}
							//
							// fmt.Println("Application", string(jsonApplication))
							if application.ScalabilitySpecifications.Data().Replicas >= domain.MaxNumberOfReplicas && domain.IsAtMaxCPULimit(*application.ContainerSpecifications.Data().CPULimit) && domain.IsAtMaxMemoryLimit(*application.ContainerSpecifications.Data().MemoryLimit) {
								fmt.Println("Manual application", application.Name, "has reached the maximum number of replicas, and is at max cpu and memory limit, email not sent to", application.AdministratorEmail)
								success, err := scheduler.scalabilityNotificationService.SendApplicationCannotBeScaledUp(
									application.AdministratorEmail,
									application.Name,
									application.Namespace.Name,
									compareActualUsageToAcceptedPercentageResults,
									application.ContainerSpecifications.Data(),
									application.ScalabilitySpecifications.Data(),
								)
								if err != nil {
									fmt.Println("error when try to send application wont be able to be scaled up notification mail during NotifyApplicationManualScalingRecommendationScheduler :", err.Error())
									done <- true
									return
								}
								if success {
									fmt.Println("Manual application", application.Name, "cannot won't be able to be scaled up, email sent to", application.AdministratorEmail)

									lastNotifiedDatetimeByApplication[application.ID] = time.Now()
								} else {
									fmt.Println("Manual application", application.Name, "cannot won't be able to be scaled up, but email not sent to", application.AdministratorEmail)
								}
								done <- true
								return
							}

							success, err := scheduler.scalabilityNotificationService.SendApplicationScalabilityRecommandationMail(
								application.AdministratorEmail,
								application.Name,
								application.Namespace.Name,
								compareActualUsageToAcceptedPercentageResults,
								application.ContainerSpecifications.Data(),
								application.ScalabilitySpecifications.Data(),
							)
							if err != nil {
								fmt.Println("error when try to send application scalability recommandation mail during NotifyApplicationManualScalingRecommendationScheduler :", err.Error())
								done <- true
								return
							}
							if success {
								fmt.Println("Manual application", application.Name, "has exceeded accepted percentages, email sent to", application.AdministratorEmail)

								lastNotifiedDatetimeByApplication[application.ID] = time.Now()
							} else {
								fmt.Println("Manual application", application.Name, "has exceeded accepted percentages, but email not sent to", application.AdministratorEmail)
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

func getManualScaleAppSchedulerRepeatInterval() (int, error) {
	var repeatInterval int
	SchedulerRecommendApplicationScalingInSeconds := os.Getenv("SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS")
	if SchedulerRecommendApplicationScalingInSeconds == "" {
		fmt.Println("SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS is not set")
		return 0, fmt.Errorf("SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS is not set")
	}
	repeatInterval, err := strconv.Atoi(SchedulerRecommendApplicationScalingInSeconds)
	if err != nil {
		fmt.Println("Error when convert SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS to int during NotifyApplicationManualScalingRecommendationScheduler :", err.Error())
		return 0, fmt.Errorf("error when convert SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS to int during NotifyApplicationManualScalingRecommendationScheduler : %s", err.Error())
	}
	return repeatInterval, nil
}
