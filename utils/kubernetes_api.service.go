package utils

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

// ConnectToKubernetesAPI Connect to Kubernetes API and return the clientset
func ConnectToKubernetesAPI() (*kubernetes.Clientset, error) {
	kubeconfig := os.Getenv("KUBECONFIG_PATH")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, nil
}

func GetKubernetesServices(clientset *kubernetes.Clientset, namespace string) (*v1.ServiceList, error) {
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return services, nil
}
