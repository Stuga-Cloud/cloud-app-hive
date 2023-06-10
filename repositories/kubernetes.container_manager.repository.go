package repositories

import (
	"bytes"
	customErrors "cloud-app-hive/controllers/errors"
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/utils"
	"k8s.io/apimachinery/pkg/api/resource"

	//"cloud-app-hive/utils"
	"context"
	"fmt"
	"io"
	//"k8s.io/apimachinery/pkg/api/resource"
	"os"
	"strings"

	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v13 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type KubernetesContainerManagerRepository struct{}

// connectToKubernetesAPIMetrics Connect to Kubernetes API and return the clientset
func (containerManager KubernetesContainerManagerRepository) connectToKubernetesAPIMetrics() (*versioned.Clientset, error) {
	kubeconfig := os.Getenv("KUBECONFIG_PATH")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, customErrors.ContainerManagerConnectionError{Message: fmt.Sprintf("Error while connecting to Kubernetes API : %s", err.Error())}
	}

	clientSet, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, customErrors.ContainerManagerConnectionError{Message: fmt.Sprintf("Error while connecting to Kubernetes API : %s", err.Error())}
	}

	return clientSet, nil
}

func (containerManager KubernetesContainerManagerRepository) GetApplicationMetrics(
	application commands.GetApplicationMetrics,
) ([]domain.ApplicationMetrics, error) {
	applicationNamespace := application.Namespace
	applicationName := application.Name

	metricsClientset, err := containerManager.connectToKubernetesAPIMetrics()
	if err != nil {
		return nil, err
	}

	clientSet, err := containerManager.connectToKubernetesAPI()
	if err != nil {
		return nil, err
	}

	metrics, err := metricsClientset.MetricsV1beta1().PodMetricses(applicationNamespace).List(
		context.Background(), metav1.ListOptions{},
	)
	if err != nil {
		return nil, customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Error while connecting to Kubernetes API : %s", err.Error()),
		}
	}

	deploymentName := fmt.Sprintf("%s-deployment", applicationName)

	// Get the Deployment object
	deployment, err := clientSet.AppsV1().Deployments(applicationNamespace).Get(
		context.Background(), deploymentName, metav1.GetOptions{},
	)
	if err != nil {
		return nil, customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Error while retrieving Deployment object : %s", err.Error()),
		}
	}

	var applicationMetrics []domain.ApplicationMetrics
	for _, metric := range metrics.Items {
		if strings.HasPrefix(metric.Name, deploymentName) {
			for _, container := range metric.Containers {
				var currentApplicationMetrics domain.ApplicationMetrics
				currentApplicationMetrics.PodName = metric.Name
				currentApplicationMetrics.Name = container.Name
				currentApplicationMetrics.CPUUsage = container.Usage.Cpu().String()
				currentApplicationMetrics.MemoryUsage = container.Usage.Memory().String()
				currentApplicationMetrics.EphemeralStorageUsage = container.Usage.StorageEphemeral().String()
				currentApplicationMetrics.PodsUsage = container.Usage.Pods().String()

				// Get deployment limits
				// Get resource limits from the Deployment object
				for _, containerSpec := range deployment.Spec.Template.Spec.Containers {
					if containerSpec.Name == container.Name {
						if containerSpec.Resources.Limits != nil {
							currentApplicationMetrics.MaxCPUUsage = containerSpec.Resources.Limits.Cpu().String()
							currentApplicationMetrics.MaxMemoryUsage = containerSpec.Resources.Limits.Memory().String()
							currentApplicationMetrics.MaxEphemeralStorage = containerSpec.Resources.Limits.Storage().String()
						}
						break
					}
				}

				applicationMetrics = append(applicationMetrics, currentApplicationMetrics)
			}
		}
	}

	return applicationMetrics, nil
}

// connectToKubernetesAPI Connect to Kubernetes API and return the clientset
func (containerManager KubernetesContainerManagerRepository) connectToKubernetesAPI() (*kubernetes.Clientset, error) {
	kubeconfig := os.Getenv("KUBECONFIG_PATH")
	if kubeconfig == "" {
		return nil, customErrors.ContainerManagerConnectionError{Message: "KUBECONFIG_PATH environment variable is not set"}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, customErrors.ContainerManagerConnectionError{
			Message: fmt.Sprintf("Error while connecting to Kubernetes API : %s", err.Error()),
		}
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, customErrors.ContainerManagerConnectionError{
			Message: fmt.Sprintf("Error while connecting to Kubernetes API : %s", err.Error()),
		}
	}

	return clientSet, nil
}

func (containerManager KubernetesContainerManagerRepository) ApplyApplication(
	applyApplication commands.ApplyApplication,
) error {
	clientset, err := containerManager.connectToKubernetesAPI()
	if err != nil {
		return err
	}

	err = containerManager.applyNamespace(clientset, applyApplication)
	if err != nil {
		return customErrors.ContainerManagerError{
			Message: "While applying namespace - " + err.Error(),
		}
	}

	err = containerManager.applyDeployment(clientset, applyApplication)
	if err != nil {
		return customErrors.ContainerManagerError{
			Message: "While applying deployment - " + err.Error(),
		}
	}

	err = containerManager.applyService(clientset, applyApplication)
	if err != nil {
		return customErrors.ContainerManagerError{
			Message: "While applying service - " + err.Error(),
		}
	}

	err = containerManager.applyIngress(clientset, applyApplication)
	if err != nil {
		return customErrors.ContainerManagerError{
			Message: "While applying ingress - " + err.Error(),
		}
	}

	return nil
}

func (containerManager KubernetesContainerManagerRepository) applyNamespace(
	clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication,
) error {
	namespace := deployApplication.Namespace
	_, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		_, err = clientset.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			return customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while creating namespace : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	}

	fmt.Println("Namespace created successfully : ", namespace)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) applyDeployment(clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	applicationImage := deployApplication.Image
	// applicationPort := deployApplication.Port
	applicationEnvironmentVariables := make([]v1.EnvVar, 0)
	for _, environmentVariable := range deployApplication.EnvironmentVariables {
		applicationEnvironmentVariables = append(applicationEnvironmentVariables, v1.EnvVar{
			Name:  environmentVariable.Name,
			Value: environmentVariable.Val,
		})
	}
	var replicas int32
	if deployApplication.ApplicationType == domain.SingleInstance {
		replicas = 1
	} else {
		replicas = deployApplication.ScalabilitySpecifications.Replicas
	}
	// TODO : Add CPU, Memory and Storage limits
	rawCpuLimit := fmt.Sprintf("%d%s", deployApplication.ContainerSpecifications.CPULimit.Val, deployApplication.ContainerSpecifications.CPULimit.Unit)
	rawMemoryLimit := fmt.Sprintf("%d%s", deployApplication.ContainerSpecifications.MemoryLimit.Val, deployApplication.ContainerSpecifications.MemoryLimit.Unit)
	cpuLimit := resource.MustParse(utils.ConvertReadableHumanValueAndUnitToK8sResource(rawCpuLimit))
	memoryLimit := resource.MustParse(utils.ConvertReadableHumanValueAndUnitToK8sResource(rawMemoryLimit))

	deploymentName := fmt.Sprintf("%s-deployment", applicationName)

	deployment := &v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: applicationNamespace,
		},
		Spec: v12.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  applicationName,
							Image: applicationImage,
							// Declare environment variables from deployApplication.EnvironmentVariables array
							Env: applicationEnvironmentVariables,
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:    cpuLimit,
									v1.ResourceMemory: memoryLimit,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := clientset.AppsV1().Deployments(applicationNamespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err == nil {
		_, err = clientset.AppsV1().Deployments(applicationNamespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
		if err != nil {
			return customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while updating deployment : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	} else {
		_, err = clientset.AppsV1().Deployments(applicationNamespace).Create(context.Background(), deployment, metav1.CreateOptions{})
		if err != nil {
			return customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while creating deployment : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	}

	fmt.Println("Deployment created successfully : " + deploymentName + " in namespace " + applicationNamespace)
	return nil
}

// Add secrets to the application
func (containerManager KubernetesContainerManagerRepository) applySecrets(clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationSecrets := deployApplication.Secrets

	for _, secret := range applicationSecrets {
		secretName := secret.Name
		secretVal := secret.Val
		secretKey := secret.Name

		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: applicationNamespace,
			},
			StringData: map[string]string{
				secretKey: secretVal,
			},
		}

		_, err := clientset.CoreV1().Secrets(applicationNamespace).Get(context.Background(), secretName, metav1.GetOptions{})
		if err == nil {
			_, err = clientset.CoreV1().Secrets(applicationNamespace).Update(context.Background(), secret, metav1.UpdateOptions{})
			if err != nil {
				return customErrors.ContainerManagerApplicationDeploymentError{
					Message:         fmt.Sprintf("Error while updating secret : %s", err.Error()),
					ApplicationName: deployApplication.Name,
					Namespace:       deployApplication.Namespace,
					Image:           deployApplication.Image,
				}
			}
		} else {
			_, err = clientset.CoreV1().Secrets(applicationNamespace).Create(context.Background(), secret, metav1.CreateOptions{})
			if err != nil {
				return customErrors.ContainerManagerApplicationDeploymentError{
					Message:         fmt.Sprintf("Error while creating secret : %s", err.Error()),
					ApplicationName: deployApplication.Name,
					Namespace:       deployApplication.Namespace,
					Image:           deployApplication.Image,
				}
			}
		}

		fmt.Println("Secret created successfully : " + secretName + " in namespace " + applicationNamespace)
	}

	return nil
}
func FrenchReadableResourceUnitToKubernetesCPUUnit(resourceUnit domain.ContainerMemoryLimitUnit) string {
	switch resourceUnit {
	case domain.MB:
		return "M"
	case domain.GB:
		return "G"
	case domain.TB:
		return "T"
	case domain.KB:
		return "K"
	default:
		return "K"
	}
}

func FrenchReadableResourceUnitToKubernetesMemoryUnit(resourceUnit domain.ContainerMemoryLimitUnit) string {
	switch resourceUnit {
	case domain.MB:
		return "Mi"
	case domain.GB:
		return "Gi"
	case domain.TB:
		return "Ti"
	case domain.KB:
		return "Ki"
	default:
		return "Ki"
	}
}

func (containerManager KubernetesContainerManagerRepository) applyService(clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	applicationPort := deployApplication.Port
	servicePort := 80
	serviceName := fmt.Sprintf("%s-service", applicationName)
	deploymentName := fmt.Sprintf("%s-deployment", applicationName)

	var serviceType v1.ServiceType
	if deployApplication.ApplicationType == domain.SingleInstance {
		serviceType = v1.ServiceTypeClusterIP
	} else if deployApplication.ApplicationType == domain.LoadBalanced {
		serviceType = v1.ServiceTypeLoadBalancer // TODO - not working with type load balancer
	} else {
		return customErrors.ContainerManagerApplicationDeploymentError{
			Message:         fmt.Sprintf("Error while creating service : %s", "Application type not supported"),
			ApplicationName: deployApplication.Name,
			Namespace:       deployApplication.Namespace,
			Image:           deployApplication.Image,
		}
	}

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: applicationNamespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app": deploymentName,
			},
			Ports: []v1.ServicePort{
				{
					Protocol:   v1.ProtocolTCP,
					Port:       int32(servicePort),
					TargetPort: intstr.FromInt(int(applicationPort)),
				},
			},
			Type: serviceType,
		},
	}

	_, err := clientset.CoreV1().Services(applicationNamespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err == nil {
		_, err = clientset.CoreV1().Services(applicationNamespace).Update(context.Background(), service, metav1.UpdateOptions{})
		if err != nil {
			return customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while updating service : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	} else {
		_, err = clientset.CoreV1().Services(applicationNamespace).Create(context.Background(), service, metav1.CreateOptions{})
		if err != nil {
			return customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while creating service : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	}

	fmt.Println("Service created successfully : " + serviceName + " in namespace " + applicationNamespace)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) applyIngress(clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	ingressName := fmt.Sprintf("%s-ingress", applicationName)
	serviceName := fmt.Sprintf("%s-service", applicationName)
	servicePort := 80
	domainName := os.Getenv("DOMAIN_NAME")
	ingress := v13.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: applicationNamespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: v13.IngressSpec{
			Rules: []v13.IngressRule{
				{
					Host: fmt.Sprintf("%s.%s.%s", applicationNamespace, applicationName, domainName),
					IngressRuleValue: v13.IngressRuleValue{
						HTTP: &v13.HTTPIngressRuleValue{
							Paths: []v13.HTTPIngressPath{
								{
									Path:     "/",
									PathType: func() *v13.PathType { p := v13.PathTypePrefix; return &p }(),
									Backend: v13.IngressBackend{
										Service: &v13.IngressServiceBackend{
											Name: serviceName,
											Port: v13.ServiceBackendPort{
												Number: int32(servicePort),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := clientset.NetworkingV1().Ingresses(applicationNamespace).Get(context.Background(), ingressName, metav1.GetOptions{})
	if err == nil {
		_, err = clientset.NetworkingV1().Ingresses(applicationNamespace).Update(context.Background(), &ingress, metav1.UpdateOptions{})
		if err != nil {
			return customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while updating ingress : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	} else {
		_, err = clientset.NetworkingV1().Ingresses(applicationNamespace).Create(context.Background(), &ingress, metav1.CreateOptions{})
		if err != nil {
			return customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while creating ingress : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	}

	fmt.Println("Ingress created successfully : " + ingressName + " in namespace " + applicationNamespace)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) GetApplicationLogs(deployApplication commands.GetApplicationLogs) ([]domain.ApplicationLogs, error) {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	deploymentName := fmt.Sprintf("%s-deployment", applicationName)

	clientset, err := containerManager.connectToKubernetesAPI()
	if err != nil {
		return nil, customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Connecting to Kubernetes API while getting application logs failed : %s", err.Error()),
		}
	}
	podList, err := clientset.CoreV1().Pods(applicationNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deploymentName),
	})
	if err != nil {
		return nil, customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Getting pods while getting application logs failed : %s", err.Error()),
		}
	}
	if len(podList.Items) == 0 {
		return nil, customErrors.ContainerManagerError{
			Message: fmt.Sprintf("No pod found for application %s in namespace %s", applicationName, applicationNamespace),
		}
	}

	logs := make([]domain.ApplicationLogs, 0)
	podLogOptions := v1.PodLogOptions{
		// TODO : Add options in body
	}

	for _, pod := range podList.Items {
		request := clientset.CoreV1().Pods(applicationNamespace).GetLogs(pod.Name, &podLogOptions)
		podLogs, err := request.Stream(context.Background())
		if err != nil {
			return nil, customErrors.ContainerManagerError{
				Message: fmt.Sprintf("Opening stream to pod %s while getting application logs failed : %s", pod.Name, err.Error()),
			}
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		if _, err = io.Copy(buf, podLogs); err != nil {
			return nil, customErrors.ContainerManagerError{
				Message: fmt.Sprintf(
					"Reading stream from pod %s while getting application logs failed : %s", pod.Name, err.Error(),
				),
			}
		}

		logs = append(logs, domain.ApplicationLogs{
			PodName: pod.Name,
			Logs:    buf.String(),
		})
	}

	return logs, nil
}

func (containerManager KubernetesContainerManagerRepository) UnapplyApplication(
	unapplyApplication commands.UnapplyApplication,
) error {
	applicationNamespace := unapplyApplication.Namespace
	applicationName := unapplyApplication.Name

	clientset, err := containerManager.connectToKubernetesAPI()
	if err != nil {
		return customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Connecting to Kubernetes API while unapplying application failed : %s", err.Error()),
		}
	}

	if err = containerManager.deleteIngress(clientset, unapplyApplication); err != nil {
		// Redeploy application if ingress deletion failed ?
		return customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Deleting ingress while unapplying application failed : %s", err.Error()),
		}
	}

	if err = containerManager.deleteService(clientset, unapplyApplication); err != nil {
		return customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Deleting service while unapplying application failed : %s", err.Error()),
		}
	}

	if err = containerManager.deleteDeployment(clientset, unapplyApplication); err != nil {
		return customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Deleting deployment while unapplying application failed : %s", err.Error()),
		}
	}

	if err = containerManager.deletePods(clientset, unapplyApplication); err != nil {
		return customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Deleting pods while unapplying application failed : %s", err.Error()),
		}
	}

	fmt.Println("Application deleted successfully : " + applicationName + " in namespace " + applicationNamespace)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) deleteIngress(clientset *kubernetes.Clientset, deployApplication commands.UnapplyApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	ingressName := fmt.Sprintf("%s-ingress", applicationName)
	err := clientset.NetworkingV1().Ingresses(applicationNamespace).Delete(context.Background(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		return customErrors.ContainerManagerApplicationRemoveError{
			Message:         fmt.Sprintf("Error deleting ingress : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
		}
	}
	fmt.Println("Ingress deleted successfully : " + applicationName)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) deleteService(clientset *kubernetes.Clientset, deployApplication commands.UnapplyApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	serviceName := fmt.Sprintf("%s-service", applicationName)
	err := clientset.CoreV1().Services(applicationNamespace).Delete(context.Background(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		return customErrors.ContainerManagerApplicationRemoveError{
			Message:         fmt.Sprintf("Error deleting service : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
		}
	}
	fmt.Println("Service deleted successfully : " + applicationName)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) deleteDeployment(clientset *kubernetes.Clientset, deployApplication commands.UnapplyApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	deploymentName := fmt.Sprintf("%s-deployment", applicationName)
	err := clientset.AppsV1().Deployments(applicationNamespace).Delete(context.Background(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		return customErrors.ContainerManagerApplicationRemoveError{
			Message:         fmt.Sprintf("Error deleting deployment : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
		}
	}
	fmt.Println("Deployment deleted successfully : " + applicationName)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) deletePods(clientset *kubernetes.Clientset, deployApplication commands.UnapplyApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	deploymentName := fmt.Sprintf("%s-deployment", applicationName)
	podList, err := clientset.CoreV1().Pods(applicationNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deploymentName),
	})
	if err != nil {
		return customErrors.ContainerManagerApplicationRemoveError{
			Message:         fmt.Sprintf("Error getting pods : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
		}
	}
	for _, pod := range podList.Items {
		err = clientset.CoreV1().Pods(applicationNamespace).Delete(context.Background(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			return customErrors.ContainerManagerApplicationRemoveError{
				Message:         fmt.Sprintf("Error deleting pod : %s", err.Error()),
				ApplicationName: applicationName,
				Namespace:       applicationNamespace,
			}
		}
		fmt.Println("Pod deleted successfully : " + pod.Name)
	}
	return nil
}

//func (containerManager KubernetesContainerManagerRepository) deleteNamespace(clientset *kubernetes.Clientset, deployApplication commands.UnapplyApplication) error {
//	applicationNamespace := deployApplication.NamespaceID
//	err := clientset.CoreV1().Namespaces().Delete(context.Background(), applicationNamespace, metav1.DeleteOptions{})
//	if err != nil {
//		return customErrors.ContainerManagerApplicationRemoveError{
//			Message:         fmt.Sprintf("Error deleting namespace : %s", err.Error()),
//			ApplicationName: deployApplication.Name,
//			Namespace:       deployApplication.NamespaceID,
//		}
//	}
//	fmt.Println("NamespaceID deleted successfully : " + applicationNamespace)
//	return nil
//}

func (containerManager KubernetesContainerManagerRepository) GetApplicationStatus(deployApplication commands.GetApplicationStatus) (*domain.ApplicationStatus, error) {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	deploymentName := fmt.Sprintf("%s-deployment", applicationName)

	clientset, err := containerManager.connectToKubernetesAPI()
	if err != nil {
		return nil, customErrors.ContainerManagerApplicationInformationError{
			Message:         fmt.Sprintf("Connecting to Kubernetes API gettings application status failed : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
			Type:            "ConnectToKubernetesAPI",
		}
	}

	deployment, err := clientset.AppsV1().Deployments(applicationNamespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, customErrors.ContainerManagerApplicationInformationError{
			Message:         fmt.Sprintf("Getting deployment failed : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
			Type:            "Deployment",
		}
	}

	serviceName := fmt.Sprintf("%s-service", applicationName)
	service, err := clientset.CoreV1().Services(applicationNamespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, customErrors.ContainerManagerApplicationInformationError{
			Message:         fmt.Sprintf("Getting service failed : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
			Type:            "Service",
		}
	}
	fmt.Printf("Service status : %s", service.Status.String())

	ingressName := fmt.Sprintf("%s-ingress", applicationName)

	ingress, err := clientset.NetworkingV1().Ingresses(applicationNamespace).Get(
		context.Background(), ingressName, metav1.GetOptions{},
	)
	if err != nil {
		return nil, customErrors.ContainerManagerApplicationInformationError{
			Message:         fmt.Sprintf("Getting ingress failed : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
			Type:            "Ingress",
		}
	}
	fmt.Println("Ingress status : " + ingress.Status.String())

	var deploymentConditions []domain.DeploymentCondition
	for _, condition := range deployment.Status.Conditions {
		deploymentConditions = append(deploymentConditions, domain.DeploymentCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastUpdateTime:     condition.LastUpdateTime.String(),
			LastTransitionTime: condition.LastTransitionTime.String(),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}

	applicationStatus := domain.ApplicationStatus{
		DeploymentName:      deployment.Name,
		StatusInString:      deployment.Status.String(),
		Replicas:            deployment.Status.Replicas,
		AvailableReplicas:   deployment.Status.AvailableReplicas,
		UnavailableReplicas: deployment.Status.UnavailableReplicas,
		ReadyReplicas:       deployment.Status.ReadyReplicas,
		DesiredReplicas:     deployment.Status.Replicas,
		CurrentReplicas:     deployment.Status.Replicas,
		UpdatedReplicas:     deployment.Status.UpdatedReplicas,
		DeploymentCondition: deploymentConditions,
		ServiceStatus:       domain.ServiceStatus{
			// TODO
			// Name:              service.Name,
			// ClusterIP:         service.Status.LoadBalancer
			// Type:              string(service.Spec.Type),
			// StatusInString:    service.Status.String(),
		},
		IngressStatus: domain.IngressStatus{
			// TODO
		},
	}

	return &applicationStatus, nil
}
