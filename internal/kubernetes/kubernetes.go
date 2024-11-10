package kubernetes

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesClient struct {
	clientset *kubernetes.Clientset
}

func NewKubernetesClient() (*KubernetesClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &KubernetesClient{
		clientset: clientset,
	}, nil
}

func (kc *KubernetesClient) GetNamespaces(ctx context.Context) ([]string, error) {
	namespaces, err := kc.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	var ns []string
	for _, namespace := range namespaces.Items {
		ns = append(ns, namespace.Name)
	}

	return ns, nil
}

type Deployment struct {
	Name         string
	Replicas     int32
	TrueReplicas int32
	Labels       map[string]string
	Containers   []Container
}

type Container struct {
	Name  string
	Image string
	Tag   string
}

func (kc *KubernetesClient) GetDeploymentsFromNamespace(ctx context.Context, namespace string) ([]Deployment, error) {
	deployments, err := kc.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	var ds []Deployment
	for _, d := range deployments.Items {
		var containers []Container
		for _, c := range d.Spec.Template.Spec.Containers {
			splitImageTag := strings.Split(c.Image, ":")
			containers = append(containers, Container{
				Name:  c.Name,
				Image: splitImageTag[0],
				Tag:   splitImageTag[1],
			})
		}

		ds = append(ds, Deployment{
			Name:         d.Name,
			Replicas:     *d.Spec.Replicas,
			TrueReplicas: d.Status.Replicas,
			Labels:       d.Labels,
			Containers:   containers,
		})
	}

	return ds, nil
}
