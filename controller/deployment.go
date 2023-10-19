package controller

import (
	"context"
	"fmt"

	"github.com/sjian_mstr/cluster-management/kubernetes"
	"github.com/sjian_mstr/cluster-management/tools/clientcmd"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListDeployments(env string, namespace string) ([]*appsv1.Deployment, error) {
	// Create a Kubernetes clientset using the provided kubeconfig.
	config, err := clientcmd.BuildConfigFromFlags("", kube[env].Configpath)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		// os.Exit(1)
	}

	// Create a Kubernetes clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating clientset: %v\n", err)
		// os.Exit(1)
	}

	// List deployments.
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %v\n", err)
		// os.Exit(1)
	}

	// Convert the slice of appsv1.Deployment to a slice of *appsv1.Deployment.
	var deploymentPointers []*appsv1.Deployment
	for i := range deployments.Items {
		deploymentPointers = append(deploymentPointers, &deployments.Items[i])
	}

	return deploymentPointers, nil
}

func DescribeDeployment(env string, deploymentName string, namespace string) {

}
