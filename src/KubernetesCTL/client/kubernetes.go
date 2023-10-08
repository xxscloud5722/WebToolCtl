package client

import (
	"context"
	appv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"path"
	"strings"
)

type KClient struct {
	Name   string
	client kubernetes.Clientset
}

// NewKClient 创建客户端实例
func NewKClient(configPath string) (*KClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}
	restClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	var ext = path.Ext(configPath)
	var temp = path.Base(configPath)
	var temps = strings.Split(temp, "\\")
	temp = temps[len(temps)-1]
	temps = strings.Split(temp, "/")
	var fileName = temps[len(temps)-1]
	return &KClient{client: *restClient, Name: strings.TrimSuffix(fileName, ext)}, nil
}

func (k *KClient) Namespaces() ([]v1.Namespace, error) {
	list, err := k.client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *KClient) Namespace(name string) (*v1.Namespace, error) {
	result, err := k.client.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (k *KClient) Deployments(namespace string) ([]appv1.Deployment, error) {
	list, err := k.client.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (k *KClient) StatefulSets(namespace string) ([]appv1.StatefulSet, error) {
	list, err := k.client.AppsV1().StatefulSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *KClient) DaemonSets(namespace string) ([]appv1.DaemonSet, error) {
	list, err := k.client.AppsV1().DaemonSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *KClient) Jobs(namespace string) ([]batchv1.Job, error) {
	list, err := k.client.BatchV1().Jobs(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *KClient) CronJobs(namespace string) ([]batchv1.CronJob, error) {
	list, err := k.client.BatchV1().CronJobs(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *KClient) Configmaps(namespace string) ([]v1.ConfigMap, error) {
	list, err := k.client.CoreV1().ConfigMaps(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *KClient) Secrets(namespace string) ([]v1.Secret, error) {
	list, err := k.client.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *KClient) Services(namespace string) ([]v1.Service, error) {
	list, err := k.client.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *KClient) Ingress(namespace string) ([]networkingv1.Ingress, error) {
	list, err := k.client.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (k *KClient) PersistentVolumeClaims(namespace string) ([]v1.PersistentVolumeClaim, error) {
	list, err := k.client.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *KClient) PersistentVolumes() ([]v1.PersistentVolume, error) {
	list, err := k.client.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
