package schedulers

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/services"
	"cloud-app-hive/use_cases"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type NotifyAdminOnClusterExceededUsageScheduler struct {
	getClusterMetricsUseCase use_cases.GetClusterMetricsUseCase
	emailService             services.EmailService
}

func (scheduler NotifyAdminOnClusterExceededUsageScheduler) Launch() {
	fmt.Println("Starting 'NotifyAdminOnClusterExceededUsageScheduler' scheduler...")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic during NotifyAdminOnClusterExceededUsageScheduler:", r)
				// relaunch scheduler
				scheduler.Launch()
			}
		}()
		repeatInterval := getNotifyAdminOnClusterExceededUsageRepeatInterval()
		ticker := time.NewTicker(time.Duration(repeatInterval) * time.Second)

		var lastNotificationSentToAdmin time.Time

		for {
			select {
			case <-ticker.C:
				clusterMetrics, err := scheduler.getClusterMetricsUseCase.Execute()
				if err != nil {
					fmt.Println("error when try to get cluster metrics during NotifyAdminOnClusterExceededUsageScheduler :", err.Error())
					panic("Error when try to get cluster metrics during NotifyAdminOnClusterExceededUsageScheduler : " + err.Error())
				}
				if clusterMetrics == nil {
					fmt.Println("No cluster metrics found")
					continue
				}

				notifyAdminWhenClusterNodesUsageIsAbovePercentageStr := os.Getenv("NOTIFY_ADMIN_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE")
				if notifyAdminWhenClusterNodesUsageIsAbovePercentageStr == "" {
					fmt.Println("NOTIFY_ADMIN_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE is not set")
					panic("NOTIFY_ADMIN_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE is not set")
				}
				notifyAdminWhenPercentageOfNodesExceedsUsageStr := os.Getenv("NOTIFY_ADMIN_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE")
				if notifyAdminWhenPercentageOfNodesExceedsUsageStr == "" {
					fmt.Println("NOTIFY_ADMIN_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE is not set")
					panic("NOTIFY_ADMIN_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE is not set")
				}

				notifyAdminWhenClusterNodesUsageIsAbovePercentage, err := strconv.ParseFloat(notifyAdminWhenClusterNodesUsageIsAbovePercentageStr, 64)
				if err != nil {
					fmt.Println("Error when convert NOTIFY_ADMIN_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE to float64")
					panic(fmt.Sprintf("Error when convert NOTIFY_ADMIN_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE to float64 during NotifyAdminOnClusterExceededUsageScheduler : %s", err.Error()))
				}
				notifyAdminWhenPercentageOfNodesExceedsUsage, err := strconv.ParseFloat(notifyAdminWhenPercentageOfNodesExceedsUsageStr, 64)
				if err != nil {
					fmt.Println("Error when convert NOTIFY_ADMIN_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE to float64")
					panic(fmt.Sprintf("Error when convert NOTIFY_ADMIN_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE to float64 during NotifyAdminOnClusterExceededUsageScheduler : %s", err.Error()))
				}

				if domain.DoesPartOfNodesExceedCPUOrMemoryUsage(
					notifyAdminWhenClusterNodesUsageIsAbovePercentage,
					notifyAdminWhenPercentageOfNodesExceedsUsage,
					clusterMetrics.NodesComputedUsages,
				) {
					fmt.Println("Cluster is exceeding its limits !")
					clusterStateJSON, err := json.Marshal(clusterMetrics)
					if err != nil {
						fmt.Println("Error while marshalling cluster state: ", err)
						continue
					}
					fmt.Println("Cluster state: ", string(clusterStateJSON))
					if time.Since(lastNotificationSentToAdmin).Hours() > 4 {
						sendClusterExceededUsageEmailToAdmin(scheduler.emailService, clusterMetrics)
						lastNotificationSentToAdmin = time.Now()
					} else {
						fmt.Println("Email already sent to admin less than 8 hours ago")
					}
				}
			}
		}
	}()
}

// NodesMetrics
// NodesCapacities
func findCorrespondingNodeMetricsAndCapacities(
	nodeName string,
	clusterMetrics *domain.ClusterMetrics,
) (domain.NodeMetrics, domain.NodeCapacities) {
	for _, nodeMetrics := range clusterMetrics.NodesMetrics {
		if nodeMetrics.Name == nodeName {
			for _, nodeCapacities := range clusterMetrics.NodesCapacities {
				if nodeCapacities.Name == nodeName {
					return nodeMetrics, nodeCapacities
				}
			}
		}
	}
	panic(fmt.Sprintf("No corresponding node metrics and capacities found for node %s", nodeName))
}

func sendClusterExceededUsageEmailToAdmin(emailService services.EmailService, clusterMetrics *domain.ClusterMetrics) {
	subject := "Cluster is exceeding its limits"
	body := "Cluster is exceeding its limits. Here are cluster metrics : \n"
	for _, nodeComputedUsage := range clusterMetrics.NodesComputedUsages {
		nodeMetrics, nodeCapacities := findCorrespondingNodeMetricsAndCapacities(
			nodeComputedUsage.Name,
			clusterMetrics,
		)

		body += fmt.Sprintf(
			"Node %s\n",
			nodeComputedUsage.Name,
		)
		body += fmt.Sprintf(
			`
			- CPU = using %2.f%% of total CPU - absolute value : %v / capacity : %v CPU\n
			- Memory = using %2.f%% of total memory - absolute value : %v / capacity : %v\n
			`,
			nodeComputedUsage.CPUUsageInPercentage,
			nodeMetrics.ReadableCPUUsage,
			nodeCapacities.ReadableCPU,
			nodeComputedUsage.MemoryUsageInPercentage,
			nodeMetrics.ReadableMemoryUsage,
			nodeCapacities.ReadableMemory,
		)
	}

	htmlBody := fmt.Sprintf(
		`
		<p>Cluster is exceeding its limits. Here are cluster metrics : </p>
		<ul>
		`,
	)
	for _, nodeComputedUsage := range clusterMetrics.NodesComputedUsages {
		nodeMetrics, nodeCapacities := findCorrespondingNodeMetricsAndCapacities(
			nodeComputedUsage.Name,
			clusterMetrics,
		)

		htmlBody += fmt.Sprintf(
			`
			<li>Node %s</li>
			<ul>
				<li>CPU = using %2.f%% of total CPU - absolute value : %v / capacity : %v</li>
				<li>Memory = using %2.f%% of total memory - absolute value : %v / capacity : %v</li>
			</ul>
			`,
			nodeComputedUsage.Name,
			nodeComputedUsage.CPUUsageInPercentage,
			nodeMetrics.ReadableCPUUsage,
			nodeCapacities.ReadableCPU,
			nodeComputedUsage.MemoryUsageInPercentage,
			nodeMetrics.ReadableMemoryUsage,
			nodeCapacities.ReadableMemory,
		)
	}
	htmlBody += "</ul>"

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		panic("ADMIN_EMAIL is not set")
	}
	copyCarbonCopyEmails := []string{}
	copyCarbonEmails := os.Getenv("ADMIN_EMAIL_COPY_CARBON")
	if copyCarbonEmails == "" {
		panic("ADMIN_EMAIL_COPY_CARBON is not set")
	}
	copyCarbonCopyEmails = strings.Split(copyCarbonEmails, ",")

	fmt.Printf("Sending %s email to admin %s and copy carbon copy emails %v\n", subject, adminEmail, copyCarbonCopyEmails)
	err := emailService.Send(
		adminEmail,
		subject,
		body,
		htmlBody,
		copyCarbonCopyEmails,
	)
	if err != nil {
		fmt.Println("error when try to send application scalability notification mail :", err.Error())
		panic("Error when try to send application scalability notification mail during NotifyApplicationScalingRecommendationScheduler : " + err.Error())
	}
}

func getNotifyAdminOnClusterExceededUsageRepeatInterval() int {
	const defaultRepeatInterval = 1
	var repeatInterval int
	SchedulerRecommendApplicationScalingInSeconds := os.Getenv("SCHEDULER_NOTIFY_ADMIN_ON_CLUSTER_EXCEEDED_USAGE_IN_SECONDS")
	if SchedulerRecommendApplicationScalingInSeconds == "" {
		fmt.Println("SCHEDULER_NOTIFY_ADMIN_ON_CLUSTER_EXCEEDED_USAGE_IN_SECONDS is not set, using default value")
		repeatInterval = defaultRepeatInterval
	}
	repeatInterval, err := strconv.Atoi(SchedulerRecommendApplicationScalingInSeconds)
	if err != nil {
		fmt.Println("Error when convert SCHEDULER_NOTIFY_ADMIN_ON_CLUSTER_EXCEEDED_USAGE_IN_SECONDS to int")
		panic(fmt.Sprintf("Error when convert SCHEDULER_NOTIFY_ADMIN_ON_CLUSTER_EXCEEDED_USAGE_IN_SECONDS to int during NotifyApplicationScalingRecommendationScheduler : %s", err.Error()))
	}
	return repeatInterval
}
