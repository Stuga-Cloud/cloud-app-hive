package repositories

import (
	"bytes"
	customErrors "cloud-app-hive/controllers/errors"
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"encoding/base64"
	"encoding/json"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"

	"context"
	"fmt"
	"io"

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
	kubeconfigContent := os.Getenv("KUBECONFIG_CONTENT")
	if kubeconfigContent == "" {
		return nil, &customErrors.ContainerManagerConnectionError{Message: "KUBECONFIG_CONTENT environment variable is not set"}
	}

	decodedContent, err := base64.StdEncoding.DecodeString(kubeconfigContent)
	if err != nil {
		// fmt.Println("Error while decoding KUBECONFIG_CONTENT, maybe its not stored as base64 encoded string : ", err.Error())
		decodedContent = []byte(kubeconfigContent)
		// return nil, &customErrors.ContainerManagerConnectionError{
		// 	Message: fmt.Sprintf("Error while decoding KUBECONFIG_CONTENT : %s", err.Error()),
		// }
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(decodedContent)
	if err != nil {
		return nil, &customErrors.ContainerManagerConnectionError{
			Message: fmt.Sprintf("Error while connecting to Kubernetes Metrics API : %s", err.Error()),
		}
	}

	clientSet, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, &customErrors.ContainerManagerConnectionError{Message: fmt.Sprintf("Error while connecting to Kubernetes Metrics API : %s", err.Error())}
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
		return nil, &customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Error while connecting to Kubernetes API : %s", err.Error()),
		}
	}

	deploymentName := fmt.Sprintf("%s-deployment", applicationName)

	// Get the Deployment object
	deployment, err := clientSet.AppsV1().Deployments(applicationNamespace).Get(
		context.Background(), deploymentName, metav1.GetOptions{},
	)
	if err != nil {
		return nil, &customErrors.ContainerManagerError{
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
	kubeconfigContent := os.Getenv("KUBECONFIG_CONTENT")
	if kubeconfigContent == "" {
		return nil, &customErrors.ContainerManagerConnectionError{Message: "KUBECONFIG_CONTENT environment variable is not set"}
	}

	decodedContent, err := base64.StdEncoding.DecodeString(kubeconfigContent)
	if err != nil {
		// fmt.Println("Error while decoding KUBECONFIG_CONTENT, maybe its not stored as base64 encoded string : ", err.Error())
		decodedContent = []byte(kubeconfigContent)
		// return nil, &customErrors.ContainerManagerConnectionError{
		// 	Message: fmt.Sprintf("Error while decoding KUBECONFIG_CONTENT : %s", err.Error()),
		// }
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(decodedContent)
	if err != nil {
		return nil, &customErrors.ContainerManagerConnectionError{
			Message: fmt.Sprintf("Error while connecting to Kubernetes API : %s", err.Error()),
		}
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, &customErrors.ContainerManagerConnectionError{
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
		return &customErrors.ContainerManagerError{
			Message: "While applying namespace - " + err.Error(),
		}
	}

	secretOriginalKeyWithConvertedK8sKey, err := containerManager.applySecrets(clientset, applyApplication)
	if err != nil {
		return &customErrors.ContainerManagerError{
			Message: "While applying secrets - " + err.Error(),
		}
	}

	if usesPrivateRegistry(applyApplication.Registry) {
		err = containerManager.applyPrivateRegistrySecret(clientset, applyApplication)
		if err != nil {
			return &customErrors.ContainerManagerError{
				Message: "While applying private registry secret - " + err.Error(),
			}
		}
	}

	err = containerManager.applyDeployment(clientset, applyApplication, secretOriginalKeyWithConvertedK8sKey)
	if err != nil {
		return &customErrors.ContainerManagerError{
			Message: "While applying deployment - " + err.Error(),
		}
	}

	err = containerManager.applyService(clientset, applyApplication)
	if err != nil {
		return &customErrors.ContainerManagerError{
			Message: "While applying service - " + err.Error(),
		}
	}

	err = containerManager.applyIngress(clientset, applyApplication)
	if err != nil {
		return &customErrors.ContainerManagerError{
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
			return &customErrors.ContainerManagerApplicationDeploymentError{
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

// DockerRegistrySecretData represents the data to store in the Docker registry Secret
type DockerRegistrySecretData struct {
	Auths map[string]struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Auth     string `json:"auth"`
	} `json:"auths"`
}

func usesPrivateRegistry(registry domain.ImageRegistry) bool {
	return registry == domain.PrivateRegistry
}

func (containerManager KubernetesContainerManagerRepository) applyPrivateRegistrySecret(
	clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication,
) error {
	privateRegistryUrl := os.Getenv("PRIVATE_HARBOR_REGISTRY_URL")
	privateRegistryUsername := os.Getenv("PRIVATE_HARBOR_REGISTRY_USERNAME")
	privateRegistryPassword := os.Getenv("PRIVATE_HARBOR_REGISTRY_PASSWORD")
	secretName := fmt.Sprintf("%s-private-registry-secret", deployApplication.Name)
	encodedAuth := base64.URLEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s:%s", privateRegistryUsername, privateRegistryPassword),
	))
	// s.Auth = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.Username, s.Password)))
	// data := k8sDockerRegistrySecretData{
	// 	Auths: map[string]k8sDockerRegistrySecret{
	// 		registry: s,
	// 	},
	// }
	// authBytes, _ := json.Marshal(data)
	//     // delete this line
	// // reStr := base64.URLEncoding.EncodeToString(authBytes)
	// return authBytes

	// Create the data for the Secret
	secretData := DockerRegistrySecretData{
		Auths: map[string]struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Auth     string `json:"auth"`
		}{
			privateRegistryUrl: {
				Username: privateRegistryUsername,
				Password: privateRegistryPassword,
				Auth:     encodedAuth,
			},
		},
	}

	// Convert the Secret data to JSON
	secretDataJSON, err := json.Marshal(secretData)
	if err != nil {
		return err
	}

	applicationNamespace := deployApplication.Namespace
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: applicationNamespace,
		},
		Type: v1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			".dockerconfigjson": []byte(secretDataJSON),
		},
	}

	_, err = clientset.CoreV1().Secrets(applicationNamespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err == nil {
		_, err = clientset.CoreV1().Secrets(applicationNamespace).Update(context.Background(), secret, metav1.UpdateOptions{})
		if err != nil {
			return &customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while updating private registry secret : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	} else {
		_, err = clientset.CoreV1().Secrets(applicationNamespace).Create(context.Background(), secret, metav1.CreateOptions{})
		if err != nil {
			return &customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while creating private registry secret : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	}

	fmt.Println("Private registry secret created successfully : " + secretName + " in namespace " + applicationNamespace)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) applyDeployment(clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication, secretOriginalKeyWithConvertedK8sKey map[string]string) error {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	applicationImage := deployApplication.Image
	if usesPrivateRegistry(deployApplication.Registry) {
		privateRegistryUrl := os.Getenv("PRIVATE_HARBOR_REGISTRY_URL")
		applicationImage = fmt.Sprintf("%s/%s", privateRegistryUrl, deployApplication.Image)
	}
	// applicationPort := deployApplication.Port
	applicationEnvironmentVariables := make([]v1.EnvVar, 0)
	for _, environmentVariable := range deployApplication.EnvironmentVariables {
		applicationEnvironmentVariables = append(applicationEnvironmentVariables, v1.EnvVar{
			Name:  environmentVariable.Name,
			Value: environmentVariable.Val,
		})
	}

	secretName := fmt.Sprintf("%s-secrets", applicationName)
	// Add secret keys to environment variables
	for secretOriginalKey, convertedK8sKey := range secretOriginalKeyWithConvertedK8sKey {
		applicationEnvironmentVariables = append(applicationEnvironmentVariables, v1.EnvVar{
			Name: secretOriginalKey,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					Key: convertedK8sKey,
					LocalObjectReference: v1.LocalObjectReference{
						Name: secretName,
					},
				},
			},
		})
	}

	var replicas int32
	if deployApplication.ApplicationType == domain.SingleInstance {
		replicas = 1
	} else {
		if deployApplication.ScalabilitySpecifications.Replicas > domain.MaxNumberOfReplicas {
			return &customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while creating deployment : %s", "Replicas must be less than or equal to "+fmt.Sprintf("%d", domain.MaxNumberOfReplicas)),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
		replicas = deployApplication.ScalabilitySpecifications.Replicas
	}
	// json, _ := json.Marshal(deployApplication)
	// fmt.Println("Deploying application with replicas : ", string(json))
	rawCpuLimit := fmt.Sprintf("%d%s", deployApplication.ContainerSpecifications.CPULimit.Val, deployApplication.ContainerSpecifications.CPULimit.Unit)
	rawMemoryLimit := fmt.Sprintf("%d%s", deployApplication.ContainerSpecifications.MemoryLimit.Val, deployApplication.ContainerSpecifications.MemoryLimit.Unit)
	cpuLimit := resource.MustParse(domain.ConvertReadableHumanValueAndUnitToK8sResource(rawCpuLimit))
	memoryLimit := resource.MustParse(domain.ConvertReadableHumanValueAndUnitToK8sResource(rawMemoryLimit))

	deploymentName := fmt.Sprintf("%s-deployment", applicationName)

	runtimeClassName := os.Getenv("RUNTIME_CLASS_NAME")
	if runtimeClassName == "" {
		runtimeClassName = "gvisor"
	}
	deployment := &v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: applicationNamespace,
			Annotations: map[string]string{
				"app.kubernetes.io/name":            applicationName,
				"app.kubernetes.io/managedBy":       "cloud-app-hive",
				"kubectl.kubernetes.io/restartedAt": time.Now().Format(time.RFC3339),
			},
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
					RuntimeClassName: &runtimeClassName,
					Containers: []v1.Container{
						{
							Name:  applicationName,
							Image: applicationImage,
							Env:   applicationEnvironmentVariables,
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:    cpuLimit,
									v1.ResourceMemory: memoryLimit,
								},
							},
							// LivenessProbe: &v1.Probe{
							// 	ProbeHandler: v1.ProbeHandler{
							// 		Exec: &v1.ExecAction{
							// 			Command: []string{
							// 				"cat",
							// 				"/tmp/healthy",
							// 			},
							// 		},
							// 		HTTPGet: &v1.HTTPGetAction{
							// 			Path: "/health",
							// 			Port: intstr.FromInt(int(deployApplication.Port)),
							// 		},
							// 	},
							// 	InitialDelaySeconds: 10,
							// 	PeriodSeconds:       10,
							// 	FailureThreshold:    3,
							// },
						},
					},
				},
			},
		},
	}

	if usesPrivateRegistry(deployApplication.Registry) {
		fmt.Println("Using private registry in deployment")
		deployment.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{
			{
				Name: fmt.Sprintf("%s-private-registry-secret", deployApplication.Name),
			},
		}
	}

	_, err := clientset.AppsV1().Deployments(applicationNamespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err == nil {
		_, err = clientset.AppsV1().Deployments(applicationNamespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
		if err != nil {
			return &customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while updating deployment : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	} else {
		_, err = clientset.AppsV1().Deployments(applicationNamespace).Create(context.Background(), deployment, metav1.CreateOptions{})
		if err != nil {
			return &customErrors.ContainerManagerApplicationDeploymentError{
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
func (containerManager KubernetesContainerManagerRepository) applySecrets(clientset *kubernetes.Clientset, deployApplication commands.ApplyApplication) (map[string]string, error) {
	applicationNamespace := deployApplication.Namespace
	applicationSecrets := deployApplication.Secrets

	// While applying secrets -
	//Error while creating secrets : Secret "A_SECRET_ENVIRONMENT_VARIABLE" is invalid:
	// => metadata.name: Invalid value: "A_SECRET_ENVIRONMENT_VARIABLE":
	// 		a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.',
	//		and must start and end with an alphanumeric character (e.g. 'example.com', regex used for validation is
	//		'[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*') (application a-second-basic-api in namespace
	//		my-big-namespace with image williamqch/basic-api:latest failed to deploy at 2023-06-16T16:20:36Z)
	secretName := fmt.Sprintf("%s-secrets", deployApplication.Name)

	stringData := make(map[string]string)
	secretOriginalKeyWithConvertedK8sKey := make(map[string]string)
	for _, secret := range applicationSecrets {
		secretKey := strings.ToLower(secret.Name)
		secretOriginalKeyWithConvertedK8sKey[secret.Name] = secretKey
		secretVal := secret.Val

		stringData[secretKey] = secretVal
	}
	secrets := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: applicationNamespace,
		},
		StringData: stringData,
		Type:       v1.SecretTypeOpaque,
	}

	_, err := clientset.CoreV1().Secrets(applicationNamespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err == nil {
		_, err = clientset.CoreV1().Secrets(applicationNamespace).Update(context.Background(), secrets, metav1.UpdateOptions{})
		if err != nil {
			return nil, &customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while updating secrets : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	} else {
		_, err = clientset.CoreV1().Secrets(applicationNamespace).Create(context.Background(), secrets, metav1.CreateOptions{})
		if err != nil {
			return nil, &customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while creating secrets : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	}

	fmt.Println("Secret created successfully : " + secretName + " in namespace " + applicationNamespace)

	return secretOriginalKeyWithConvertedK8sKey, nil
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
		serviceType = v1.ServiceTypeLoadBalancer
	} else {
		return &customErrors.ContainerManagerApplicationDeploymentError{
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
			return &customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while updating service : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	} else {
		_, err = clientset.CoreV1().Services(applicationNamespace).Create(context.Background(), service, metav1.CreateOptions{})
		if err != nil {
			return &customErrors.ContainerManagerApplicationDeploymentError{
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
				"cert-manager.io/cluster-issuer":             "letsencrypt",
			},
		},
		Spec: v13.IngressSpec{
			IngressClassName: func() *string { s := "nginx"; return &s }(),
			TLS: []v13.IngressTLS{
				{
					Hosts: []string{
						fmt.Sprintf("%s.%s.%s", applicationName, applicationNamespace, domainName),
					},
					SecretName: "letsencrypt-account-key",
				},
			},
			Rules: []v13.IngressRule{
				{
					Host: fmt.Sprintf("%s.%s.%s", applicationName, applicationNamespace, domainName),
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
			return &customErrors.ContainerManagerApplicationDeploymentError{
				Message:         fmt.Sprintf("Error while updating ingress : %s", err.Error()),
				ApplicationName: deployApplication.Name,
				Namespace:       deployApplication.Namespace,
				Image:           deployApplication.Image,
			}
		}
	} else {
		_, err = clientset.NetworkingV1().Ingresses(applicationNamespace).Create(context.Background(), &ingress, metav1.CreateOptions{})
		if err != nil {
			return &customErrors.ContainerManagerApplicationDeploymentError{
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
		return nil, &customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Connecting to Kubernetes API while getting application logs failed : %s", err.Error()),
		}
	}
	podList, err := clientset.CoreV1().Pods(applicationNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deploymentName),
	})
	if err != nil {
		return nil, &customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Getting pods while getting application logs failed : %s", err.Error()),
		}
	}
	if len(podList.Items) == 0 {
		return nil, &customErrors.ContainerManagerError{
			Message: fmt.Sprintf("No pod found for application %s in namespace %s", applicationName, applicationNamespace),
		}
	}

	logs := make([]domain.ApplicationLogs, 0)
	podLogOptions := v1.PodLogOptions{
		// Follow: true,
		// Timestamps: true,
	}

	for _, pod := range podList.Items {
		request := clientset.CoreV1().Pods(applicationNamespace).GetLogs(pod.Name, &podLogOptions)
		podLogs, err := request.Stream(context.Background())
		if err != nil {
			return nil, &customErrors.ContainerManagerError{
				Message: fmt.Sprintf("Opening stream to pod %s while getting application logs failed : %s", pod.Name, err.Error()),
			}
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		if _, err = io.Copy(buf, podLogs); err != nil {
			return nil, &customErrors.ContainerManagerError{
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
		return &customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Connecting to Kubernetes API while unapplying application failed : %s", err.Error()),
		}
	}

	if err = containerManager.deleteIngress(clientset, unapplyApplication); err != nil {
		// TODO: Redeploy application if ingress deletion failed ?
		return &customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Deleting ingress while unapplying application failed : %s", err.Error()),
		}
	}

	if err = containerManager.deleteService(clientset, unapplyApplication); err != nil {
		return &customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Deleting service while unapplying application failed : %s", err.Error()),
		}
	}

	if err = containerManager.deleteDeployment(clientset, unapplyApplication); err != nil {
		return &customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Deleting deployment while unapplying application failed : %s", err.Error()),
		}
	}

	if err = containerManager.deletePods(clientset, unapplyApplication); err != nil {
		return &customErrors.ContainerManagerError{
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
		return &customErrors.ContainerManagerApplicationRemoveError{
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
		return &customErrors.ContainerManagerApplicationRemoveError{
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
		return &customErrors.ContainerManagerApplicationRemoveError{
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
		return &customErrors.ContainerManagerApplicationRemoveError{
			Message:         fmt.Sprintf("Error getting pods : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
		}
	}
	for _, pod := range podList.Items {
		err = clientset.CoreV1().Pods(applicationNamespace).Delete(context.Background(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			return &customErrors.ContainerManagerApplicationRemoveError{
				Message:         fmt.Sprintf("Error deleting pod : %s", err.Error()),
				ApplicationName: applicationName,
				Namespace:       applicationNamespace,
			}
		}
		fmt.Println("Pod deleted successfully : " + pod.Name)
	}
	return nil
}

func (containerManager KubernetesContainerManagerRepository) DeleteNamespace(namespace string) error {
	clientset, err := containerManager.connectToKubernetesAPI()
	if err != nil {
		return &customErrors.ContainerManagerError{
			Message: fmt.Sprintf("Connecting to Kubernetes API while deleting namespace failed : %s", err.Error()),
		}
	}

	// Namespace can be not created yet
	_, err = clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return &customErrors.ContainerManagerNamespaceRemoveError{
			Message:   fmt.Sprintf("Error getting namespace while deleting namespace : %s", err.Error()),
			Namespace: namespace,
		}
	}

	err = clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		return &customErrors.ContainerManagerNamespaceRemoveError{
			Message:   fmt.Sprintf("Error deleting namespace %s : %s", namespace, err.Error()),
			Namespace: namespace,
		}
	}
	fmt.Println("NamespaceID deleted successfully : " + namespace)
	return nil
}

func (containerManager KubernetesContainerManagerRepository) GetApplicationStatus(deployApplication commands.GetApplicationStatus) (*domain.ApplicationStatus, error) {
	applicationNamespace := deployApplication.Namespace
	applicationName := deployApplication.Name
	deploymentName := fmt.Sprintf("%s-deployment", applicationName)

	clientset, err := containerManager.connectToKubernetesAPI()
	if err != nil {
		return nil, &customErrors.ContainerManagerApplicationInformationError{
			Message:         fmt.Sprintf("Connecting to Kubernetes API gettings application status failed : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
			Type:            "ConnectToKubernetesAPI",
		}
	}

	deployment, err := clientset.AppsV1().Deployments(applicationNamespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, &customErrors.ContainerManagerApplicationInformationError{
			Message:         fmt.Sprintf("Getting deployment failed : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
			Type:            "Deployment",
		}
	}

	pods, err := clientset.CoreV1().Pods(applicationNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deploymentName),
	})
	if err != nil {
		return nil, &customErrors.ContainerManagerApplicationInformationError{
			Message:         fmt.Sprintf("Getting pods failed : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
			Type:            "Pods",
		}
	}
	podList := domain.ConvertPods(pods)
	podList.Items = domain.ComputeHumanizedPodStatus(&podList.Items)
	fmt.Println(podList.Items[0].HumanizedStatus)

	serviceName := fmt.Sprintf("%s-service", applicationName)
	service, err := clientset.CoreV1().Services(applicationNamespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, &customErrors.ContainerManagerApplicationInformationError{
			Message:         fmt.Sprintf("Getting service failed : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
			Type:            "Service",
		}
	}

	ingressName := fmt.Sprintf("%s-ingress", applicationName)

	ingress, err := clientset.NetworkingV1().Ingresses(applicationNamespace).Get(
		context.Background(), ingressName, metav1.GetOptions{},
	)
	if err != nil {
		return nil, &customErrors.ContainerManagerApplicationInformationError{
			Message:         fmt.Sprintf("Getting ingress failed : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
			Type:            "Ingress",
		}
	}

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

	// Order the deploymentConditions slice if there are different last update times
	// isDifferentLastUpdateTimes := false
	// for i := 0; i < len(deploymentConditions)-1; i++ {
	// 	if deploymentConditions[i].LastUpdateTime != deploymentConditions[i+1].LastUpdateTime {
	// 		isDifferentLastUpdateTimes = true
	// 		break
	// 	}
	// }
	// if len(deploymentConditions) > 1 && isDifferentLastUpdateTimes {
	// 	sort.Slice(deploymentConditions, func(i, j int) bool {
	// 		if deploymentConditions[i].LastUpdateTime == deploymentConditions[j].LastUpdateTime {
	// 			return deploymentConditions[i].LastTransitionTime > deploymentConditions[j].LastTransitionTime
	// 		}
	// 		return deploymentConditions[i].LastUpdateTime > deploymentConditions[j].LastUpdateTime
	// 	})
	// }

	// REVERSE the deploymentConditions slice
	//for i := len(deploymentConditions)/2 - 1; i >= 0; i-- {
	//	opp := len(deploymentConditions) - 1 - i
	//	deploymentConditions[i], deploymentConditions[opp] = deploymentConditions[opp], deploymentConditions[i]
	//}

	applicationStatus := domain.ApplicationStatus{
		Name:                deployment.Name,
		StatusInString:      deployment.Status.String(),
		Replicas:            deployment.Status.Replicas,
		AvailableReplicas:   deployment.Status.AvailableReplicas,
		UnavailableReplicas: deployment.Status.UnavailableReplicas,
		ReadyReplicas:       deployment.Status.ReadyReplicas,
		DesiredReplicas:     deployment.Status.Replicas,
		CurrentReplicas:     deployment.Status.Replicas,
		UpdatedReplicas:     deployment.Status.UpdatedReplicas,
		DeploymentCondition: deploymentConditions,
		PodList:             podList,
		ServiceStatus: domain.ServiceStatus{
			Name: serviceName,
			Type: string(service.Spec.Type),
			IP:   service.Spec.ClusterIP,
			Port: service.Spec.Ports[0].Port,
		},
		IngressStatus: domain.IngressStatus{
			Name: ingress.Name,
		},
	}

	computedStatus, humanizedStatus, err := applicationStatus.ComputeApplicationStatus()
	if err != nil {
		return nil, &customErrors.ContainerManagerApplicationInformationError{
			Message:         fmt.Sprintf("Computing application status failed : %s", err.Error()),
			ApplicationName: applicationName,
			Namespace:       applicationNamespace,
			Type:            "ComputeApplicationStatus",
		}
	}
	applicationStatus.ComputedApplicationStatus = computedStatus
	applicationStatus.HumanizedStatus = humanizedStatus

	return &applicationStatus, nil
}

// GetKubeClusterState returns the state of the cluster, CPU usage, memory usage, pods, services, deployments, namespaces
func (containerManager KubernetesContainerManagerRepository) GetClusterMetrics() (*domain.ClusterMetrics, error) {
	metricsClientset, err := containerManager.connectToKubernetesAPIMetrics()
	if err != nil {
		return nil, err
	}

	nodeMetricsList, err := metricsClientset.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeMetrics []domain.NodeMetrics
	for _, nodeMetric := range nodeMetricsList.Items {
		nodeMetrics = append(nodeMetrics, domain.NodeMetrics{
			Name:                  nodeMetric.Name,
			CPUUsage:              nodeMetric.Usage.Cpu().String(),
			MemoryUsage:           nodeMetric.Usage.Memory().String(),
			StorageUsage:          nodeMetric.Usage.Storage().String(),
			EphemeralStorageUsage: nodeMetric.Usage.StorageEphemeral().String(),
			Pods:                  nodeMetric.Usage.Pods().String(),

			ReadableCPUUsage: domain.ConvertK8sResourceToReadableHumanValueAndUnit(nodeMetric.Usage.Cpu().String()),
			ReadableMemoryUsage: domain.ConvertK8sResourceToReadableHumanValueAndUnit(
				nodeMetric.Usage.Memory().String(),
			),
			ReadableStorageUsage: domain.ConvertK8sResourceToReadableHumanValueAndUnit(
				nodeMetric.Usage.Storage().String(),
			),
			ReadableEphemeralStorageUsage: domain.ConvertK8sResourceToReadableHumanValueAndUnit(
				nodeMetric.Usage.StorageEphemeral().String(),
			),
		})
	}

	// Also get maximum CPU and memory capacity of the cluster
	clientset, err := containerManager.connectToKubernetesAPI()
	if err != nil {
		return nil, err
	}

	nodeList, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeCapacities []domain.NodeCapacities
	for _, node := range nodeList.Items {
		nodeCapacities = append(nodeCapacities, domain.NodeCapacities{
			Name:                     node.Name,
			CPULimit:                 node.Status.Capacity.Cpu().String(),
			MemoryLimit:              node.Status.Capacity.Memory().String(),
			StorageLimit:             node.Status.Capacity.Storage().String(),
			EphemeralStorageLimit:    node.Status.Capacity.StorageEphemeral().String(),
			ReadableCPU:              domain.ConvertK8sResourceToReadableHumanValueAndUnit(node.Status.Capacity.Cpu().String()),
			ReadableMemory:           domain.ConvertK8sResourceToReadableHumanValueAndUnit(node.Status.Capacity.Memory().String()),
			ReadableStorage:          domain.ConvertK8sResourceToReadableHumanValueAndUnit(node.Status.Capacity.Storage().String()),
			ReadableEphemeralStorage: domain.ConvertK8sResourceToReadableHumanValueAndUnit(node.Status.Capacity.StorageEphemeral().String()),
		})
	}

	nodesComputedUsages, err := domain.ComputeNodesUsagesFromMetricsAndCapacities(
		nodeMetrics, nodeCapacities,
	)
	if err != nil {
		fmt.Println("Error while computing nodes usages : " + err.Error())
		return nil, err
	}

	clusterMetrics := domain.ClusterMetrics{
		NodesMetrics:        nodeMetrics,
		NodesCapacities:     nodeCapacities,
		NodesComputedUsages: nodesComputedUsages,
	}

	return &clusterMetrics, nil
}
