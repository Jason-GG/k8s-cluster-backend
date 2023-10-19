package controller

import (
	"context"

	"github.com/sjian_mstr/cluster-management/kubernetes"
	"github.com/sjian_mstr/cluster-management/tools/clientcmd"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListNamespacesTest() ([]*v1.Namespace, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kube["dev"].Configpath)
	if err != nil {
		return nil, err
	}

	// Create a Kubernetes clientset using the provided kubeconfig.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// List Namespaces.
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// Convert the slice of v1.Namespace to a slice of *v1.Namespace.
	var namespacePointers []*v1.Namespace
	for i := range namespaces.Items {
		namespacePointers = append(namespacePointers, &namespaces.Items[i])
	}

	return namespacePointers, nil

}

func ListNamespaces(cluster string) ([]*v1.Namespace, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kube[cluster].Configpath)
	if err != nil {
		return nil, err
	}

	// Create a Kubernetes clientset using the provided kubeconfig.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// List Namespaces.
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// Convert the slice of v1.Namespace to a slice of *v1.Namespace.
	var namespacePointers []*v1.Namespace
	for i := range namespaces.Items {
		namespacePointers = append(namespacePointers, &namespaces.Items[i])
	}

	return namespacePointers, nil
}
