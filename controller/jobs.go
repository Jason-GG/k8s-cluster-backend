package controller

import (
	"context"
	"fmt"

	"github.com/sjian_mstr/cluster-management/kubernetes"
	"github.com/sjian_mstr/cluster-management/tools/clientcmd"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListJodsInfo(env string, namespace string) (*batchv1.JobList, error) {
	// Path to kubeconfig file (default location is $HOME/.kube/config)

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

	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return jobs, nil

}
