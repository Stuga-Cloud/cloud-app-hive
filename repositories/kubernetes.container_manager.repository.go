package repositories

import (
	"bytes"
	"cloud-app-hive/domain"
	"context"
	"fmt"
	"io"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v13 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

func (containerManager KubernetesContainerManagerRepository) GetApplicationMetrics(application domain.GetApplicationMetrics) ([]domain.ApplicationMetrics, error) {
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

func (containerManager KubernetesContainerManagerRepository) DeployApplication(deployApplication domain.DeployApplication) error {
	clientset := containerManager.connectToKubernetesAPI()

	err := containerManager.createNamespaceIfNotExists(clientset, deployApplication)
	if err != nil {
		panic(err.Error())
	}

	err = containerManager.createDeployment(clientset, deployApplication)
	if err != nil {
		panic(err.Error())
	}

	err = containerManager.createService(clientset, deployApplication)
	if err != nil {
		panic(err.Error())
	}

	err = containerManager.createIngress(clientset, deployApplication)
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func (containerManager KubernetesContainerManagerRepository) createNamespaceIfNotExists(clientset *kubernetes.Clientset, deployApplication domain.DeployApplication) error {
	namespace := deployApplication.Namespace
	list, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	fmt.Println("Error: ", err)
	fmt.Println("Namespace: ", list)
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

func (containerManager KubernetesContainerManagerRepository) createDeployment(clientset *kubernetes.Clientset, deployApplication domain.DeployApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	applicationImage := deployApplication.Image
	var replicas int32
	if deployApplication.ApplicationType == domain.ApplicationType(0) {
		replicas = 1
	} else {
		replicas = deployApplication.ScalabilitySpecifications.Replicas
	}
	cpuLimit := deployApplication.ContainerSpecifications.CpuLimit
	memoryLimit := deployApplication.ContainerSpecifications.MemoryLimit
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
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%d%s", cpuLimit.Value, cpuLimit.Unit.String())),
									v1.ResourceMemory: resource.MustParse(fmt.Sprintf("%d%s", memoryLimit.Value, memoryLimit.Unit.String())),
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

func (containerManager KubernetesContainerManagerRepository) createService(clientset *kubernetes.Clientset, deployApplication domain.DeployApplication) error {
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
		serviceType = v1.ServiceTypeLoadBalancer
	}

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
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

func (containerManager KubernetesContainerManagerRepository) createIngress(clientset *kubernetes.Clientset, deployApplication domain.DeployApplication) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	ingressName := fmt.Sprintf("%s-ingress", applicationName)
	serviceName := fmt.Sprintf("%s-service", applicationName)
	servicePort := 80
	domainName := os.Getenv("DOMAIN_NAME")
	ingress := v13.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: ingressName,
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

func (containerManager KubernetesContainerManagerRepository) GetApplicationLogs(deployApplication domain.GetApplicationLogs) ([]domain.ApplicationLogs, error) {
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
		return nil, fmt.Errorf("no pod found for application %s in namespace %s", applicationName, applicationNamespace)
	}

	logs := make([]domain.ApplicationLogs, 0)
	podLogOptions := v1.PodLogOptions{
		// TODO
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

//func (containerManager KubernetesContainerManagerRepository) deleteDeployment(clientset *kubernetes.Clientset, deployApplication domain.DeployApplication) (string, error) {
//	applicationNamespace := deployApplication.Namespace
//	applicationName := deployApplication.PodName
//	deploymentName := fmt.Sprintf("%s-deployment", applicationName)
//	err := clientset.AppsV1().Deployments(applicationNamespace).Delete(context.Background(), deploymentName, metav1.DeleteOptions{})
//	if err != nil {
//		panic(err.Error())
//	}
//	return "Deployment deleted successfully : " + applicationName, nil
//}
