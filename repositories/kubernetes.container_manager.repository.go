package repositories

import (
	"bytes"
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"context"
	"errors"
	"fmt"
	"io"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v13 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
	"strings"
)

type KubernetesContainerManagerRepository struct{}

// connectToKubernetesAPIMetrics Connect to Kubernetes API and return the clientset
func (containerManager KubernetesContainerManagerRepository) connectToKubernetesAPIMetrics() *versioned.Clientset {
	kubeconfig := os.Getenv("KUBECONFIG_PATH")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return versioned.NewForConfigOrDie(config)
}

func (containerManager KubernetesContainerManagerRepository) GetApplicationMetrics(application commands.GetApplicationMetrics) ([]domain.ApplicationMetrics, error) {
	applicationNamespace := application.Namespace
	applicationName := application.Name

	metricsClientset := containerManager.connectToKubernetesAPIMetrics()
	metrics, err := metricsClientset.MetricsV1beta1().PodMetricses(applicationNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	deploymentName := fmt.Sprintf("%s-deployment", applicationName)
	var applicationMetrics []domain.ApplicationMetrics
	for _, metric := range metrics.Items {
		if strings.HasPrefix(metric.Name, deploymentName) {
			for _, container := range metric.Containers {
				var currentApplicationMetrics domain.ApplicationMetrics
				currentApplicationMetrics.Name = container.Name
				currentApplicationMetrics.CPUUsage = container.Usage.Cpu().String()
				currentApplicationMetrics.MemoryUsage = container.Usage.Memory().String()
				currentApplicationMetrics.StorageUsage = container.Usage.Storage().String()
				currentApplicationMetrics.EphemeralStorageUsage = container.Usage.StorageEphemeral().String()
				currentApplicationMetrics.PodsUsage = container.Usage.Pods().String()
				applicationMetrics = append(applicationMetrics, currentApplicationMetrics)
			}
		}
	}

	return applicationMetrics, nil
}

// connectToKubernetesAPI Connect to Kubernetes API and return the clientset
func (containerManager KubernetesContainerManagerRepository) connectToKubernetesAPI() *kubernetes.Clientset {
	kubeconfig := os.Getenv("KUBECONFIG_PATH")
	if kubeconfig == "" {
		panic("KUBECONFIG_PATH environment variable is not set")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return kubernetes.NewForConfigOrDie(config)
}

func (containerManager KubernetesContainerManagerRepository) ApplyApplication(applyApplication commands.ApplyApplication) error {
	clientset := containerManager.connectToKubernetesAPI()

	err := containerManager.applyNamespace(clientset, applyApplication)
	if err != nil {
		panic(err.Error())
	}

	err = containerManager.applyDeployment(clientset, applyApplication)
	if err != nil {
		panic(err.Error())
	}

	err = containerManager.applyService(clientset, applyApplication)
	if err != nil {
		panic(err.Error())
	}

	err = containerManager.applyIngress(clientset, applyApplication)
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func (containerManager KubernetesContainerManagerRepository) applyNamespace(clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication) error {
	namespace := deployApplication.Namespace
	_, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		_, err = clientset.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			panic(err.Error())
		}
	}

	fmt.Println("Namespace created successfully : ", namespace)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) applyDeployment(clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	applicationImage := deployApplication.Image
	//applicationPort := deployApplication.Port
	var applicationEnvironmentVariables = make([]v1.EnvVar, 0)
	for _, environmentVariable := range deployApplication.EnvironmentVariables {
		applicationEnvironmentVariables = append(applicationEnvironmentVariables, v1.EnvVar{
			Name:  environmentVariable.Name,
			Value: environmentVariable.Value,
		})
	}
	var replicas int32
	if deployApplication.ApplicationType == domain.ApplicationType(0) {
		replicas = 1
	} else {
		replicas = deployApplication.ScalabilitySpecifications.Replicas
	}
	//cpuLimit := deployApplication.ContainerSpecifications.CpuLimit
	//memoryLimit := deployApplication.ContainerSpecifications.MemoryLimit
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
							//Resources: v1.ResourceRequirements{ TODO : make it work and set default values
							//	Requests: v1.ResourceList{
							//		v1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%d%s", cpuLimit.Value, cpuLimit.Unit.String())),
							//		v1.ResourceMemory: resource.MustParse(fmt.Sprintf("%d%s", memoryLimit.Value, memoryLimit.Unit.String())),
							//	},
							//},
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
			panic(err.Error())
		}
	} else {
		_, err = clientset.AppsV1().Deployments(applicationNamespace).Create(context.Background(), deployment, metav1.CreateOptions{})
		if err != nil {
			panic(err.Error())
		}
	}

	fmt.Println("Deployment created successfully : " + deploymentName + " in namespace " + applicationNamespace)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) applyService(clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	applicationPort := deployApplication.Port
	servicePort := 80
	serviceName := fmt.Sprintf("%s-service", applicationName)
	deploymentName := fmt.Sprintf("%s-deployment", applicationName)

	var serviceType v1.ServiceType
	if deployApplication.ApplicationType == domain.ApplicationType(0) {
		serviceType = v1.ServiceTypeClusterIP
	} else {
		serviceType = v1.ServiceTypeLoadBalancer // TODO - not working with type load balancer
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
					TargetPort: intstr.FromInt(applicationPort),
				},
			},
			Type: serviceType,
		},
	}

	_, err := clientset.CoreV1().Services(applicationNamespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err == nil {
		_, err = clientset.CoreV1().Services(applicationNamespace).Update(context.Background(), service, metav1.UpdateOptions{})
		if err != nil {
			panic(err.Error())
		}
	} else {
		_, err = clientset.CoreV1().Services(applicationNamespace).Create(context.Background(), service, metav1.CreateOptions{})
		if err != nil {
			panic(err.Error())
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
			panic(err.Error())
		}
	} else {
		_, err = clientset.NetworkingV1().Ingresses(applicationNamespace).Create(context.Background(), &ingress, metav1.CreateOptions{})
		if err != nil {
			panic(err.Error())
		}
	}

	fmt.Println("Ingress created successfully : " + ingressName + " in namespace " + applicationNamespace)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) GetApplicationLogs(deployApplication commands.GetApplicationLogs) ([]domain.ApplicationLogs, error) {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	deploymentName := fmt.Sprintf("%s-deployment", applicationName)
	clientset := containerManager.connectToKubernetesAPI()
	podList, err := clientset.CoreV1().Pods(applicationNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deploymentName),
	})
	if err != nil {
		panic(err.Error())
	}
	if len(podList.Items) == 0 {
		return nil, errors.New("no pod found for application " + applicationName + " in namespace " + applicationNamespace)
	}

	logs := make([]domain.ApplicationLogs, 0)
	podLogOptions := v1.PodLogOptions{
		// TODO : Add options in body
	}
	for _, pod := range podList.Items {
		request := clientset.CoreV1().Pods(applicationNamespace).GetLogs(pod.Name, &podLogOptions)
		podLogs, err := request.Stream(context.Background())
		if err != nil {
			panic(err.Error())
		}
		defer podLogs.Close()
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			panic(err.Error())
		}

		logs = append(logs, domain.ApplicationLogs{
			PodName: pod.Name,
			Logs:    buf.String(),
		})
	}

	return logs, nil
}

func (containerManager KubernetesContainerManagerRepository) UnapplyApplication(unapplyApplication commands.UnapplyApplication) error {
	applicationNamespace := unapplyApplication.Namespace
	applicationName := unapplyApplication.Name

	clientset := containerManager.connectToKubernetesAPI()
	err := containerManager.deleteIngress(clientset, unapplyApplication)
	if err != nil {
		panic(err.Error())
	}

	err = containerManager.deleteService(clientset, unapplyApplication)
	if err != nil {
		panic(err.Error())
	}

	err = containerManager.deleteDeployment(clientset, unapplyApplication)
	if err != nil {
		panic(err.Error())
	}

	err = containerManager.deletePods(clientset, unapplyApplication)
	if err != nil {
		panic(err.Error())
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
		panic(err.Error())
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
		panic(err.Error())
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
		panic(err.Error())
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
		panic(err.Error())
	}
	for _, pod := range podList.Items {
		err = clientset.CoreV1().Pods(applicationNamespace).Delete(context.Background(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Pod deleted successfully : " + pod.Name)
	}
	return nil
}

//func (containerManager KubernetesContainerManagerRepository) deleteNamespace(clientset *kubernetes.Clientset, deployApplication commands.UnapplyApplication) error {
//	applicationNamespace := deployApplication.Namespace
//	err := clientset.CoreV1().Namespaces().Delete(context.Background(), applicationNamespace, metav1.DeleteOptions{})
//	if err != nil {
//		panic(err.Error())
//	}
//	fmt.Println("Namespace deleted successfully : " + applicationNamespace)
//	return nil
//}

//func (containerManager KubernetesContainerManagerRepository) GetApplicationStatus(deployApplication commands.GetApplicationStatus) (string, error) {
//	applicationNamespace := deployApplication.Namespace
//	applicationName := deployApplication.Name
//	clientset := containerManager.connectToKubernetesAPI()
//	deploymentName := fmt.Sprintf("%s-deployment", applicationName)
//	deployment, err := clientset.AppsV1().Deployments(applicationNamespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
//	if err != nil {
//		panic(err.Error())
//	}
//	if deployment.Status.ReadyReplicas == 0 {
//		return "Application " + applicationName + " is not ready", nil
//	}
//	return "Application " + applicationName + " is ready", nil
//}

//func (containerManager KubernetesContainerManagerRepository) deleteDeployment(clientset *kubernetes.Clientset, deployApplication domain.ApplyApplication) (string, error) {
//	applicationNamespace := deployApplication.Namespace
//	applicationName := deployApplication.PodName
//	deploymentName := fmt.Sprintf("%s-deployment", applicationName)
//	err := clientset.AppsV1().Deployments(applicationNamespace).Delete(context.Background(), deploymentName, metav1.DeleteOptions{})
//	if err != nil {
//		panic(err.Error())
//	}
//	return "Deployment deleted successfully : " + applicationName, nil
//}
